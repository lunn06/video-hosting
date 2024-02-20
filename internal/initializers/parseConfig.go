package initializers

import "github.com/lunn06/video-hosting/internal/config"

var CFG config.Config

func ParseConfig() {
	CFG = config.MustLoad("configs/main.yaml")
}
