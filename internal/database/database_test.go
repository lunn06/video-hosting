package database

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lunn06/video-hosting/internal/config"
	"log"
	"testing"
)

var schemaTest = `
CREATE TABLE IF NOT EXISTS person (
	id integer,
    name text 
);`

type Person struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

func TestMustCreate(t *testing.T) {
	cfg := config.MustLoad("../../configs/main.yaml")
	dbConnArg := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Address,
		cfg.Database.Name,
	)
	db, err := sqlx.Connect("pgx", dbConnArg)
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(schemaTest)
	defer db.MustExec("DROP TABLE person")
	tx := db.MustBegin()
	check := "checkQueryInsert"
	tx.MustExec("INSERT INTO person VALUES ($1, $2)", 1, check)
	tx.Commit()
	persons := []Person{}

	err = db.Select(&persons, "SELECT * FROM person WHERE name=$1", check)
	if err != nil {
		log.Fatal(err)
	}
	if check != persons[0].Name {
		log.Fatal("Error query DB: test: SELECT and INSERT")
	}
}
