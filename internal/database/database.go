package database

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lunn06/video-hosting/internal/config"
)

var schema = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER UNIQUE
);`

type User struct {
	Id int `db:"id"`
}

func MustCreate(cfg *config.Config) *sqlx.DB {
	dbConnArg := getPgAddress(cfg)

	db, err := sqlx.Connect("pgx", dbConnArg)
	if err != nil {
		log.Fatal(err)
	}

	db.MustExec(schema)

	return db
}

func getPgAddress(cfg *config.Config) string {
	// postgresql://db_user:db_password@pgadmin.dnc-check234.freemyip.com:5432/video-hosting
	return fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Address,
		cfg.Database.Name,
	)
}
