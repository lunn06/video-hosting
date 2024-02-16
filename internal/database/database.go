package database

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lunn06/video-hosting/internal/config"
	"log"
)

var schema = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER UNIQUE
);`

type User struct {
	Id int `db:"id"`
}

func MustCreate(config *config.Config) *sqlx.DB {
	// postgresql://db_user:db_password@pgadmin.dnc-check234.freemyip.com:5432/video-hosting
	dbConnArg := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Address,
		config.Database.Name,
	)

	db, err := sqlx.Connect("pgx", dbConnArg)
	if err != nil {
		log.Fatal(err)
	}

	db.MustExec(schema)

	return db
}
