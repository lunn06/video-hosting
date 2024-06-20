package services

import (
	"github.com/gin-gonic/gin"
	"github.com/lunn06/video-hosting/internal/transport/rest"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/registration", rest.Registration)
			auth.POST("/login", rest.Authentication)
			auth.POST("/refresh", rest.RefreshTokens)
			auth.GET("/ping", AuthMiddleware(), rest.Ping)
		}

		api.POST("/upload", rest.Upload)
	}

	SetupDocs(r)
	return r
}
