package initializers

import (
	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/database"
)

func ConnectToDB() {
	database.DB = database.MustCreate(config.CFG)
}
