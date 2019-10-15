package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/oauth2.v3/example/server/oauth2server"
	"gopkg.in/oauth2.v3/models"
)

func ClientRegisterHandler(c *gin.Context) {
	clients := oauth2server.GetOauth2Server().Manager.GetClients()

	if c.Request.Method == "POST" {
		name := c.PostForm("client_name")
		if name == "" {
			c.HTML(200, "client_register.html", gin.H{
				"hint":    "register failed, please input client name",
				"clients": clients,
			})
			return
		}

		description := c.PostForm("description")
		if description == "" {
			c.HTML(200, "client_register.html", gin.H{
				"hint":    "register failed, please input description",
				"clients": clients,
			})
			return
		}

		redirect_uri := c.PostForm("redirect_uri")
		if redirect_uri == "" {
			c.HTML(200, "client_register.html", gin.H{
				"hint":    "register failed, please input redirect_uri",
				"clients": clients,
			})
			return
		}

		client := models.Client{
			Name:        name,
			Description: description,
			Domain:      redirect_uri,
		}

		err := oauth2server.GetOauth2Server().Manager.RegisterClient(&client)
		if err != nil {
			c.HTML(200, "client_register.html", gin.H{
				"hint":    fmt.Sprintf("register failed, %v", err),
				"clients": clients,
			})
			return
		}

		clients = oauth2server.GetOauth2Server().Manager.GetClients()
	}

	c.HTML(200, "client_register.html", gin.H{
		"clients": clients,
	})
}
