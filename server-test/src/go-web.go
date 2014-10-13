package main

import (
	"github.com/gin-gonic/gin"
	"sym/models"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLTemplates("templates/*")
	r.GET("/", hello)
	r.GET("/json", json)
	r.POST("/user", createUser)
	r.Run(":9090")
}

func hello(c *gin.Context) {
	obj := gin.H{"title": "Main website"}
	c.HTML(200, "index.tmpl", obj)
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
		c.JSON(200, &Result{Status: true, Value: token,})
	} else {
		c.JSON(500, &Result{Status: false, Value: nil,})
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
	Value interface {}
}

type Msg struct {
	Name string `json:"user"`
	Message string
	Number int
}
