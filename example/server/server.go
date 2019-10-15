package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/oauth2.v3/example/server/oauth2server"
	"gopkg.in/oauth2.v3/example/server/services"
	"gopkg.in/oauth2.v3/example/server/web"
	"gopkg.in/oauth2.v3/server"
	"log"
)

var srv *server.Server

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	// 初始化oauth2 server
	oauth2server.Init()

	// 接口服务
	httpSvr := gin.Default()
	httpSvr.LoadHTMLGlob("static/*")

	// 服务路由定义
	test := httpSvr.Group("/oauth2/test")
	{
		// web server
		test.Any("/register", web.RegisterHandler)
		test.Any("/login", web.LoginHandler)
		test.GET("/auth", web.AuthHandler)
		test.Any("/client-register", web.ClientRegisterHandler)

		// rest api
		test.Any("/authorize", services.AuthorizeCodeHandler)
		test.POST("/access_token", services.AccessTokenHandler)
		test.POST("/refresh_token", services.RefreshTokenHandler)
	}

	log.Printf("oauth2 server is running on port 9096...")
	err := httpSvr.Run(":9096")

	if err != nil {
		log.Printf("httpSvr startup failed, %v", err)
	}
}
