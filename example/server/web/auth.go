package web

import (
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"log"
	"net/http"
)

func AuthHandler(c *gin.Context) {
	// 检查用户是否已经登录，如果没有登录，重定向到登录页面
	store, err := session.Start(nil, c.Writer, c.Request)
	if err != nil {
		c.HTML(200, "auth.html", gin.H{
			"hint": "internal error, please try again",
		})
		return
	}

	username := c.Query("username")
	if username == "" {
		c.Redirect(http.StatusFound, "/oauth2/test/login")
		return
	}

	if _, ok := store.Get(username); !ok {
		log.Printf("user %s not login", username)
		c.Redirect(http.StatusFound, "/oauth2/test/login")
		return
	}

	c.HTML(200, "auth.html", gin.H{
		"username": username,
	})
}
