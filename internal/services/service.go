package services

import (
	"github.com/gin-gonic/gin"
	"github.com/lunn06/video-hosting/internal/transport/rest"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", rest.Ping)
	r.POST("/registration", rest.Registration)
	SetupDocs(r)

	return r
}
