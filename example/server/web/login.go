package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"gopkg.in/oauth2.v3/example/server/oauth2server"
	"net/http"
)

func LoginHandler(c *gin.Context) {

	if c.Request.Method == "POST" {
		// 1. 检查用户合法性
		username := c.PostForm("username")
		if username == "" {
			c.HTML(200, "login.html", gin.H{
				"hint": "login failed, please input username",
			})
			return
		}

		password := c.PostForm("password")
		if password == "" {
			c.HTML(200, "login.html", gin.H{
				"hint": "login failed, please input password",
			})
			return
		}

		user, err := oauth2server.GetOauth2Server().Manager.GetUser(username)
		if err != nil {
			c.HTML(200, "login.html", gin.H{
				"hint": "login failed, user not existed",
			})
			return
		}

		if user.GetPassword() != password {
			c.HTML(200, "login.html", gin.H{
				"uname": username,
				"pwd":   password,
				"hint":  "login failed, wrong password",
			})
			return
		}

		// 2. 登录用户存储到session
		store, err := session.Start(nil, c.Writer, c.Request)
		if err != nil {
			c.Abort()
			return
		}

		store.Set(username, true)
		store.Save()

		// 3. 重定向到授权页面（带上userid)
		c.Redirect(http.StatusFound, fmt.Sprintf("/oauth2/test/auth?username=%s", username))

		return
	}

	c.HTML(200, "login.html", nil)
}
