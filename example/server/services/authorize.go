package services

import (
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"gopkg.in/oauth2.v3/example/server/oauth2server"
	"net/http"
	"net/url"
)

func AuthorizeCodeHandler(c *gin.Context) {
	store, err := session.Start(nil, c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}
	c.Request.Form = form

	store.Delete("ReturnUri")
	store.Save()

	err = oauth2server.GetOauth2Server().HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
}
