package services

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/oauth2.v3/example/server/oauth2server"
	"log"
)

func AccessTokenHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	err := oauth2server.GetOauth2Server().HandleTokenRequest(w, r)
	if err != nil {
		log.Printf("we're here, err %v", err)
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"error_code": -1,
			"error_msg":  err.Error,
		})
	}
}

func RefreshTokenHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	err := oauth2server.GetOauth2Server().HandleTokenRequest(w, r)
	if err != nil {
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"error_code": -1,
			"error_msg":  err.Error,
		})
	}
}
