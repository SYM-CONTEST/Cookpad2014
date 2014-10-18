package main

import (
	"github.com/SYM-CONTEST/Cookpad2014/kinennbi/crawler"
	"fmt"
	"github.com/SYM-CONTEST/Cookpad2014/kinennbi/models"
	"github.com/gin-gonic/gin"
	"github.com/mrjones/oauth"
	//"io/ioutil"
	"log"
	"net/http"
	"time"
	"strings"
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
	r.GET("/result", getResult)
	r.GET("/a/:id", getAniversary)

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
	secret, v_err := c.Request.Cookie("secret")

	var verificationCode string
	var tokenKey string
	if t_err != nil && v_err != nil {
		values := c.Request.URL.Query()
		verificationCode = values.Get("oauth_verifier")
		tokenKey = values.Get("oauth_token")
	} else {
		verificationCode = secret.Value
		tokenKey = token.Value
	}



	t_cookie := http.Cookie{
		Name:    "token",
		Value:   tokenKey,
		Expires: time.Now().AddDate(0, 1, 0),
	}

	v_cookie := http.Cookie{
		Name:    "secret",
		Value:   verificationCode,
		Expires: time.Now().AddDate(0, 1, 0),
	}

	http.SetCookie(c.Writer, &t_cookie)
	http.SetCookie(c.Writer, &v_cookie)

	c.HTML(200, "waiting.tmpl", nil)

}

func getResult(c *gin.Context) {
	token, _ := c.Request.Cookie("token")
	secret, _ := c.Request.Cookie("secret")

	oauthToken := &oauth.RequestToken{
		Token:  token.Value,
		Secret: secret.Value,
	}

	accessToken, err := consumer.AuthorizeToken(oauthToken, secret.Value)
	if err != nil {
		log.Fatal(err)
	}

	// 認証した人のcrawlerを生成
	cr := crawler.NewCrawler(accessToken.Token, accessToken.Secret)
	// 認証した人のメンションを分析してそれっぽい記念日群を抽出
	as := cr.AnalyzeAnniversary()
	// ただの確認出力なので不要
	// crawler.OutputAniversarries(as)
	// 記念日群からもっともよさげなものを選定 (今は取り急ぎ最古のもの)
	a := crawler.ChooseBestAniversary(as)
	// サイト用メッセージ1
	first := a.CreateFirstMessage()
	log.Println(first)
	// サイト用メッセージ2
	second := a.CreateSecondMessage()
	log.Println(second)

	// aniversaryをDBに保存
	aniversary := new(models.Aniversary)
	id, _ := aniversary.Create(first, second, a.Names())

	// Tweet用メッセージ
	full := a.CreateFullMessage(first, second) + " http://localhost:9090/a/" + id
	log.Println(full)
	// Tweetする時はこれで
	//	c.PostByAniv(full)

	obj := gin.H{"result": full}
	c.HTML(200, "result.tmpl", obj)
}

func getAniversary(c *gin.Context) {
	id := c.Params.ByName("id")

	aniversary := new(models.Aniversary).Get(id)

	first := aniversary.Prefix
	second := aniversary.Message
	users := strings.Split(aniversary.Users, ",")

	obj := gin.H{"first": first, "second": second, "users": users}
	c.HTML(200, "result.tmpl", obj)
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
