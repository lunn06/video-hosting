package app

import (
	"log"

	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/services"
)

func Run() {
	cfg := config.MustLoad("configs/main.yaml")

	//_ = database.MustCreateDB(cfg)

	r := services.SetupRouter()

	onlyPortAddress := ":" + cfg.HTTPServer.Port
	log.Fatal(r.Run(onlyPortAddress))
}
