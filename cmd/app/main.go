package main

import (
	"fmt"
	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/services"
)

func main() {
	cfg := config.MustLoad("configs/main.yaml")

	//_ = database.MustCreateDB(cfg)

	r := services.SetupRouter()
	err := r.Run(cfg.HTTPServer.Address)

	if err != nil {
		fmt.Println(err)
	}
}
