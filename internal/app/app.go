package app

import (
	"fmt"
	"github.com/lunn06/video-hosting/internal/services"
)

func Run() {
	r := services.SetupRouter()
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}

}
