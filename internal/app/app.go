package app

import (
	"github.com/lunn06/video-hosting/internal/initializers"
	"log"

	"github.com/lunn06/video-hosting/internal/services"
)

func init() {
	initializers.ParseConfig()
	initializers.ConnectToDB()
}

func Run() {
	r := services.SetupRouter()

	onlyPortAddress := ":" + initializers.CFG.HTTPServer.Port
	log.Fatal(r.Run(onlyPortAddress))
}
