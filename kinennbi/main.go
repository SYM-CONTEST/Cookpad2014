package main

import (
	"fmt"
	"github.com/SYM-CONTEST/Cookpad2014/kinennbi/models"
	"github.com/gin-gonic/gin"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//var tokens map[string]*oauth.RequestToken
var consumer *oauth.Consumer

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLTemplates("templates/*")
	r.Static("/static", "static")
	r.GET("/", hello)
	r.GET("/login", authTwitter)
	r.GET("/maketoken", getTwitterToken)

	// jsonレスポンスのお試し
	r.GET("/json", json)
	// ユーザ作成のテストAPI
	r.POST("/user", createUser)
	r.Run(":9090")
}

func hello(c *gin.Context) {
	obj := gin.H{"title": "Main website"}
	c.HTML(200, "index.tmpl", obj)
}

func authTwitter(c *gin.Context) {

	_, err := c.Request.Cookie("token")

	if err == nil {
		c.Redirect(302, "/maketoken")
	} else {

		//tokens = make(map[string]*oauth.RequestToken)

		consumerKey := "n567b7sH6HrPIWBZyhHM2QiaK"

		consumerSecret := "ygYGJ7aXEh2UQgLI7pOOWU5cixK6o7pDWYVY4MmvRaerJjqLwT"

		consumer = oauth.NewConsumer(
			consumerKey,
			consumerSecret,
			oauth.ServiceProvider{
				RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
				AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
				AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
			},
		)
		consumer.Debug(true)

		// コールバックURL
		tokenUrl := fmt.Sprintf("http://%s/maketoken", c.Request.Host)
		token, requestUrl, err := consumer.GetRequestTokenAndUrl(tokenUrl)
		if err != nil {
			log.Fatal(err)
		}
		// Make sure to save the token, we'll need it for AuthorizeToken()
		t := new(models.Token)
		t.Create(token.Token, token.Secret)
		//tokens[token.Token] = token

		c.Redirect(302, requestUrl)
	}
}

func getTwitterToken(c *gin.Context) {
	token, t_err := c.Request.Cookie("token")
	verifier, v_err := c.Request.Cookie("verifier")

	var verificationCode string
	var tokenKey string
	if t_err != nil && v_err != nil {
		values := c.Request.URL.Query()
		verificationCode = values.Get("oauth_verifier")
		tokenKey = values.Get("oauth_token")
	} else {
		verificationCode = verifier.Value
		tokenKey = token.Value
	}

	t := new(models.Token).Get(tokenKey)
	oauthToken := &oauth.RequestToken{
		Token:  t.Token,
		Secret: t.Secret,
	}
	accessToken, err := consumer.AuthorizeToken(oauthToken, verificationCode)
	if err != nil {
		log.Fatal(err)
	}

	response, err := consumer.Get(
		"https://api.twitter.com/1.1/statuses/home_timeline.json",
		map[string]string{"count": "1"},
		accessToken)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	t_cookie := http.Cookie{
		Name:    "token",
		Value:   tokenKey,
		Expires: time.Now().AddDate(0, 1, 0),
	}

	v_cookie := http.Cookie{
		Name:    "verifier",
		Value:   verificationCode,
		Expires: time.Now().AddDate(0, 1, 0),
	}

	http.SetCookie(c.Writer, &t_cookie)
	http.SetCookie(c.Writer, &v_cookie)

	bits, err := ioutil.ReadAll(response.Body)
	c.String(200, "The newest item in your home timeline is: "+string(bits))
}

func createUser(c *gin.Context) {
	user := new(models.User)

	var userForm struct {
		Name     string `form:"name" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	c.Bind(&userForm)

	token, result := user.Create(userForm.Name, userForm.Password)

	if result {
		c.JSON(200, &Result{Status: true, Value: token})
	} else {
		c.JSON(500, &Result{Status: false, Value: nil})
	}
}

func json(c *gin.Context) {
	msg := new(Msg)
	msg.Name = "Yuichi"
	msg.Message = "hello world"
	msg.Number = 123
	c.JSON(200, msg)
}

type Result struct {
	Status bool `json:"status"`
	Value  interface{}
}

type Msg struct {
	Name    string `json:"user"`
	Message string
	Number  int
}
