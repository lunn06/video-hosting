package app

import (
	"fmt"
	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/services"
)

func Run() {
	cfg := config.MustLoad("configs/main.yaml")

	r := services.SetupRouter()
	err := r.Run(cfg.Address)
	if err != nil {
		fmt.Println(err)
	}

}
