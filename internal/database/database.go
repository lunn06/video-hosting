package database

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

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
		return errors.New("no DB connection")
	}
	return nil
}

func InsertUser(user models.User) (uint32, error) {
	if err := checkDBConnection(); err != nil {
		slog.Error("error in insert users_token", err)
		return 0, err
	}
	tx := DB.MustBegin()
	var lastInsertIndex uint32
	err := tx.QueryRow("INSERT INTO users (email, channel_name, password) VALUES($1, $2, $3) RETURNING id", user.Email, user.ChannelName, user.Password).Scan(&lastInsertIndex)
	if err != nil {
		slog.Error("error in insert user", err)
		return 0, err
	}

	_, err = tx.Exec("INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)", lastInsertIndex, 1)
	if err != nil {
		slog.Error("error when adding user role", err)
		return 0, err
	}
	tx.Commit()
	return lastInsertIndex, err
}

func InsertVideo(video models.Video) error {
	if err := checkDBConnection(); err != nil {
		slog.Error("error in insert users_token", err)
		return err
	}

	_, err := DB.NamedExec(insertVideoRequest, video)

	if err != nil {
		return err
	}

	return err
}
func InsertToken(userId uint32, jwtToken string) (string, error) {
	if err := checkDBConnection(); err != nil {
		slog.Error("error in insert users_token", err)
		return "", err
	}
	tx := DB.MustBegin()

	var tempUUID string
	err := tx.QueryRow("INSERT INTO jwt_tokens (token) VALUES ($1) RETURNING uuid", jwtToken).Scan(&tempUUID)

	if err != nil {
		slog.Error("error in insert token", err)
		return "", err
	}

	_, err = tx.Exec("INSERT INTO users_tokens (user_id, token_uuid) VALUES ($1, $2)", userId, tempUUID)
	if err != nil {
		slog.Error("error in insert users_token", err)
		return "", err
	}
	tx.Commit()

	return tempUUID, err
}

func UpdateTokenTime(token models.JwtToken) error {
	if err := checkDBConnection(); err != nil {
		slog.Error("checkDBConnection() error = %v", err)
		return err
	}

	_, err := DB.Exec("UPDATE jwt_tokens SET creation_time=$1 WHERE uuid=$2", time.Now(), token.Uuid)
	if err != nil {
		return err
	}

	return nil
}

func GetTokenByUser(user models.User) (*models.JwtToken, error) {
	if err := checkDBConnection(); err != nil {
		slog.Error("checkDBConnection() error = %v", err)
		return nil, err
	}

	getTokenByUserRequest := `
		SELECT * FROM jwt_tokens WHERE uuid=(
		    SELECT token_uuid FROM users_tokens WHERE user_id=$1
		)
	`

	var token models.JwtToken
	err := DB.Get(&token, getTokenByUserRequest, user.Id)
	if err != nil {
		slog.Error("error on getting token from jwt_tokens table")
		return nil, err
	}

	return &token, nil
}

func GetUserByRefreshToken(token string) (*models.User, error) {
	if err := checkDBConnection(); err != nil {
		slog.Error("checkDBConnection() error = %v", err)
		return nil, err
	}

	getUserByRefreshTokenDBRequest := `
		SELECT * FROM users WHERE id=(
			SELECT user_id FROM users_tokens WHERE token_uuid=(
			    SELECT token_uuid FROM jwt_tokens WHERE token=$1
			)
		)
	`

	var user models.User
	err := DB.Get(&user, getUserByRefreshTokenDBRequest, token)
	if err != nil {
		slog.Error("GetUserByRefreshToken() error = %v, can't get tokenId from DB")
		return nil, err
	}

	return &user, nil
}

func GetToken(givenToken string) (*models.JwtToken, error) {
	if err := checkDBConnection(); err != nil {
		return nil, err
	}

	var token models.JwtToken

	err := DB.Get(&token, "SELECT * FROM jwt_tokens WHERE token=$1", givenToken)
	if err != nil {
		slog.Error(fmt.Sprintf("GetToken() error = %v, can't select", err))
		return nil, err
	}

	return &token, nil
}

func PopToken(tokenUUID string) (*models.JwtToken, error) {
	if err := checkDBConnection(); err != nil {
		return nil, err
	}

	var token models.JwtToken

	err := DB.Get(&token, "SELECT * FROM jwt_tokens WHERE uuid=$1", tokenUUID)
	if err != nil {
		slog.Error(fmt.Sprintf("PopToken() error = %v, can't select", err))
		return nil, err
	}

	_, err = DB.Exec("DELETE FROM jwt_tokens WHERE uuid=$1", tokenUUID)
	if err != nil {
		slog.Error(fmt.Sprintf("PopToken() error = %v, can't delete token", err))
		return nil, err
	}

	return &token, nil
}

func GetUser(email string) (*models.User, error) {
	if err := checkDBConnection(); err != nil {
		return nil, err
	}

	var user models.User

	err := DB.Get(&user, "SELECT * FROM users WHERE email=$1", email)

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
