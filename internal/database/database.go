package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/models"
)

var (
	DB *sqlx.DB

	insertUserRequest = `INSERT INTO users VALUES (
		DEFAULT, :email, :channel_name, :password, :registration_time
		)`

	insertVideoRequest = `INSERT INTO videos VALUES (
		DEFAULT, :title, :localization, :upload_time, :file_path, :likes_count, :views_count
		)`
)

func Init() {
	DB = MustCreate(config.CFG)

	databaseDefaults := config.MustLoadDatabaseDefaults("configs/database_defaults.yaml")

	tx := DB.MustBegin()
	for _, role := range databaseDefaults.Roles {
		tx.MustExec(
			`INSERT INTO roles VALUES (
			$1, $2, $3, $4
			) ON CONFLICT (id) DO NOTHING `,
			role.Id, role.Name, role.CanRemoveUsers, role.CanRemoveOthersVideos,
		)
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Init() error = %v", err)
	}
}

func MustCreate(cfg config.Config) *sqlx.DB {
	if DB != nil {
		return DB
	}

	dbConnArg := getPgAddress(cfg)

	db, err := sqlx.Connect("pgx", dbConnArg)
	if err != nil {
		log.Fatalf("MustCreate() Error: %v", err)
	}

	db.MustExec(models.Schema)

	return db
}

func getPgAddress(cfg config.Config) string {
	// postgresql://db_user:db_password@pgadmin.dnc-check234.freemyip.com:5432/video-hosting
	return fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Address,
		cfg.Database.Name,
	)
}

func checkDBConnection() error {
	if DB == nil {
		return errors.New("GetUser() Error: no DB connection")
	}
	return nil
}

func InsertUser(user models.User) error {
	if err := checkDBConnection(); err != nil {
		return err
	}

	_, err := DB.NamedExec(insertUserRequest, user)

	if err != nil {
		return err
	}

	return nil
}

func InsertVideo(video models.Video) error {
	if err := checkDBConnection(); err != nil {
		return err
	}

	_, err := DB.NamedExec(insertVideoRequest, video)

	if err != nil {
		return err
	}

	return nil
}

func GetUser(id uint32) (*models.User, error) {
	if err := checkDBConnection(); err != nil {
		return nil, err
	}

	var user models.User

	err := DB.Get(&user, "SELECT * FROM users WHERE id=$1", id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetVideo(id uuid.UUID) (*models.Video, error) {
	if err := checkDBConnection(); err != nil {
		return nil, err
	}

	var user models.Video

	err := DB.Get(&user, "SELECT * FROM videos WHERE uuid=$1", id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
