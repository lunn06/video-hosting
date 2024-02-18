package database

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/lunn06/video-hosting/internal/config"
	"reflect"
	"testing"
)

func TestMustCreate(t *testing.T) {
	check := 99999
	want := User{Id: check}

	cfg := config.MustLoad("../../configs/main.yaml")

	db := MustCreate(cfg)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Error on database.MustCreate(): %v", r)
		}
	}()

	tx := db.MustBegin()
	if _, err := tx.Exec("INSERT INTO users VALUES ($1)", check); err != nil {
		t.Errorf("Error on INSERT %v VALUE in users TABLE: %v", check, err)
	}

	if err := tx.Commit(); err != nil {
		t.Errorf("Error on COMMIT in users TABLE: %v", err)
	}

	user := User{}
	if err := db.Get(&user, "SELECT * FROM users WHERE id=$1", check); err != nil {
		t.Errorf("Error on db.Get: %v", err)
	}
	if !reflect.DeepEqual(user, want) {
		t.Errorf("SELECT Error: db.Get() = %v, want = %v", user, want)
	}

	if _, err := db.Exec("DELETE FROM users WHERE id=$1", check); err != nil {
		t.Errorf("DELETE Error: %v. FIX THE DB MANUALY!", err)
	}
}
