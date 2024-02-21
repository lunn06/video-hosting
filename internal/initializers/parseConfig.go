package initializers

import "github.com/lunn06/video-hosting/internal/config"

func ParseConfig() {
	config.CFG = config.MustLoad("configs/main.yaml")
}
