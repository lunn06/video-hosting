package initializers

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lunn06/video-hosting/internal/database"
)

var DB *sqlx.DB

func ConnectToDB() {
	DB = database.MustCreate(CFG)
}
