package app

import (
	"log"

	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/initializers"
	"github.com/lunn06/video-hosting/internal/services"
)

func init() {
	initializers.ParseConfig()
	initializers.ConnectToDB()
}

func Run() {
	r := services.SetupRouter()

	onlyPortAddress := ":" + config.CFG.HTTPServer.Port
	log.Fatal(r.Run(onlyPortAddress))
}
