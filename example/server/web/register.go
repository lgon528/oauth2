package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/oauth2.v3/example/server/oauth2server"
	"gopkg.in/oauth2.v3/models"
)

func RegisterHandler(c *gin.Context) {

	users := oauth2server.GetOauth2Server().Manager.GetUsers()
	if c.Request.Method == "POST" {
		// 1. 检查用户合法性
		username := c.PostForm("username")
		if username == "" {
			c.HTML(200, "register.html", gin.H{
				"hint":  "register failed, please input username",
				"users": users,
			})
			return
		}

		password := c.PostForm("password")
		if password == "" {
			c.HTML(200, "register.html", gin.H{
				"hint":  "register failed, please input password",
				"users": users,
			})
			return
		}

		// 2. 登录用户存储到session
		user := models.User{
			ID:       username,
			Password: password,
		}
		err := oauth2server.GetOauth2Server().Manager.RegisterUser(&user)
		if err != nil {
			c.HTML(200, "register.html", gin.H{
				"hint":  fmt.Sprintf("register failed, %v", err),
				"users": users,
			})
		}

		users = oauth2server.GetOauth2Server().Manager.GetUsers()
	}

	c.HTML(200, "register.html", gin.H{
		"users": users,
	})
}
