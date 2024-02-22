package initializers

import (
	"fmt"

	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/database"
)

func ConnectToDB() {
	database.DB = database.MustCreate(config.CFG)

	databaseDefaults := config.MustLoadDatabaseDefaults("configs/database_defaults.yaml")

	tx := database.DB.MustBegin()
	_, err := tx.NamedExec(
		`INSERT INTO roles VALUES (
			:id, :name, :can_remove_users, :can_remove_others_videos
		) ON CONFLICT (id) DO NOTHING `,
		databaseDefaults.Roles,
	)

	if err != nil {
		panic(fmt.Sprintf("Error in ConnectToDB: %v", err))
	}
}
