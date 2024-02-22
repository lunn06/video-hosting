package database

import (
	"reflect"
	"testing"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/models"
)

//TODO: переписать тесты, сделать тесты для всех таблицы

func TestMustCreate(t *testing.T) {
	checkID := "-2"
	checkEmail := "check@tt.com"
	checkChannelName := "test"
	checkPassword := "1234"
	want := models.User{
		Id:          checkID,
		Email:       checkEmail,
		ChannelName: checkChannelName,
		Password:    checkPassword,
	}

	cfg := config.MustLoad("../../configs/main.yaml")

	db := MustCreate(cfg)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Error on database.MustCreate(): %v", r)
		}
	}()

	tx := db.MustBegin()
	if _, err := tx.Exec("INSERT INTO users VALUES ($1, $2, $3, $4)", checkID, checkEmail, checkChannelName, checkPassword); err != nil {
		t.Errorf("Error on INSERT %v VALUE in users TABLE: %v", checkID, err)
	}

	if err := tx.Commit(); err != nil {
		t.Errorf("Error on COMMIT in users TABLE: %v", err)
	}

	user := models.User{}
	if err := db.Get(&user, "SELECT * FROM users WHERE id=$1", checkID); err != nil {
		t.Errorf("Error on db.Get: %v", err)
	}
	if !reflect.DeepEqual(user, want) {
		t.Errorf("SELECT Error: db.Get() = %v, want = %v", user, want)
	}

	if _, err := db.Exec("DELETE FROM users WHERE id=$1", checkID); err != nil {
		t.Errorf("DELETE Error: %v. FIX THE DB MANUALY!", err)
	}
}

func Test_getPgAddress(t *testing.T) {
	cfg := config.MustLoad("../../configs/example_main.yaml")

	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"getPgAdressTest",
			args{cfg},
			"postgresql://db_user:db_password@pgsql.com:5432/db_name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPgAddress(tt.args.cfg); got != tt.want {
				t.Errorf("getPgAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
