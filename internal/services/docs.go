package services

import (
	"github.com/gin-gonic/gin"
	"github.com/lunn06/video-hosting/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupDocs(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
