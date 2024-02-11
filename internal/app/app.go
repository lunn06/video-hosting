package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/services"
)

func Run() {
	cnf := config.MustLoad()
	fmt.Println(cnf)
	a := gin.New()
	_ = a
	r := services.SetupRouter()
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}

}
