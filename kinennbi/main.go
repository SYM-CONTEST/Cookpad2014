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
var tokenStore *models.Token

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

	token, err := c.Request.Cookie("accesstoken")
	tokenStore = new(models.Token)

	if err == nil || (token != nil && len(token.Value) == 0 ){
		c.Redirect(302, "/maketoken")
	} else {

		//tokens = make(map[string]*oauth.RequestToken)

		consumerKey := "WKs1pXAfWbwat1MPimOWdmoBm"

		consumerSecret := "0tvflEMDAz0yrjTEKUis3bueEVGjNhBy2tR25pNWRAZpKcAhrO"

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
	token, t_err := c.Request.Cookie("accesstoken")

	var verificationCode string
	var tokenKey string
	var accessTokenStr string
	if t_err != nil {
		values := c.Request.URL.Query()
		verificationCode = values.Get("oauth_verifier")
		tokenKey = values.Get("oauth_token")
	} else {
		accessTokenStr = token.Value
	}

	oauthToken := &oauth.RequestToken{
		Token:  tokenKey,
		Secret: verificationCode,
	}

	var accessToken *oauth.AccessToken
	var t string
	var s string
	if len(accessTokenStr) == 0 {
		accessToken, _ = consumer.AuthorizeToken(oauthToken, verificationCode)
		t = accessToken.Token
		s = accessToken.Secret
	} else {
		tmp := tokenStore.Get(accessTokenStr)
		t = tmp.Token
		s = tmp.Secret
	}

	tokenStore.Create(t, s)

	t_cookie := http.Cookie{
		Name:    "accesstoken",
		Value:   t,
		Expires: time.Now().AddDate(0, 1, 0),
	}

	http.SetCookie(c.Writer, &t_cookie)

	c.HTML(200, "waiting.tmpl", nil)

}

func getResult(c *gin.Context) {
	token, _ := c.Request.Cookie("accesstoken")

	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)

			t_cookie := http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Now().AddDate(0, 0, 0),
		}

			v_cookie := http.Cookie{
			Name:    "secret",
			Value:   "",
			Expires: time.Now().AddDate(0, 0, 0),
		}

			http.SetCookie(c.Writer, &t_cookie)
			http.SetCookie(c.Writer, &v_cookie)
			c.Redirect(301, "/")
		}
	}()

	accessToken := tokenStore.Get(token.Value)

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
	//log.Println(first)
	// サイト用メッセージ2
	second, statusId := a.CreateSecondMessage()
	//log.Println(second)

	//embed := a.GetEmbededHTML(statusId)
	embed := cr.GetOEmbed(statusId, a.Tweets)
	// aniversaryをDBに保存
	aniversary := new(models.Aniversary)
	id, _ := aniversary.Create(first, second, a.Names(), embed.Html, embed.Url)

	// Tweet用メッセージ
	var names []string = make([]string, len(a.Names()))
	for key, val := range a.Names() {
		names[key] = "@" + val
	}
	full := a.CreateFullMessage(first, second)
	log.Println(full)
	// Tweetする時はこれで
	cr.PostByAniv(full + " http://kinen.yabuchin.com/a/" + id)

	fmt.Println("full: ", full)

	obj := gin.H{"full": "." + full, "first": first, "second": second, "users": a.Names(), "embed": embed, "id": id, "embedUrl": embed.Url}
	c.HTML(200, "result.tmpl", obj)
}

func getAniversary(c *gin.Context) {
	id := c.Params.ByName("id")

	aniversary := new(models.Aniversary).Get(id)

	first := aniversary.Prefix
	second := aniversary.Message
	users := strings.Split(aniversary.Users, ",")
	embed := aniversary.Embed
	url := aniversary.URL

	var names []string = make([]string, len(users))
	for key, val := range users {
		names[key] = "@" + val
	}

	full := "." + strings.Join(names, ",") + " " + first + second

	obj := gin.H{"full": full, "first": first, "second": second, "users": users, "embed": embed, "id": id, "embedUrl": url}
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
