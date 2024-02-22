package database

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/models"
)

func init() {
	config.CFG = config.MustLoad("../../configs/main.yaml")
	DB = MustCreate(config.CFG)
}

//TODO: переписать тесты, сделать тесты для всех таблицы

var (
	checkUserID      = uuid.New().String()
	checkEmail       = "check@tt.com"
	checkChannelName = "test"
	checkPassword    = "1234"
	checkTime        = time.Now()

	testUser = models.User{
		Id:               checkUserID,
		Email:            checkEmail,
		ChannelName:      checkChannelName,
		Password:         checkPassword,
		RegistrationTime: checkTime,
	}

	checkVideoId           = uuid.New().String()
	checkTitle             = "Test Title"
	checkLocalization      = "ru_RU"
	checkUploadTime        = time.Now()
	checkFilePath          = "/some/test/dir/to/file.mp4"
	checkLikesCount   uint = 100
	checkViewsCount   uint = 10000

	testVideo = models.Video{
		Id:           checkVideoId,
		Title:        checkTitle,
		Localization: checkLocalization,
		UploadDate:   checkUploadTime,
		FilePath:     checkFilePath,
		LikesCount:   checkLikesCount,
		ViewsCount:   checkViewsCount,
	}

	insertUserRequest = `INSERT INTO users VALUES (
		:id, :email, :channel_name, :password, :registration_time
	) ON CONFLICT (id, email) DO NOTHING`
)

func TestMustCreate(t *testing.T) {
	want := testUser

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Error on database.MustCreate(): %v", r)
			return
		}
	}()

	tx := DB.MustBegin()
	if _, err := tx.Exec(
		"INSERT INTO users VALUES ($1, $2, $3, $4, $5)",
		checkUserID, checkEmail, checkChannelName, checkPassword, checkTime,
	); err != nil {
		t.Errorf("Error on INSERT %v VALUE in users TABLE: %v", checkUserID, err)
		return
	}

	if err := tx.Commit(); err != nil {
		t.Errorf("Error on COMMIT in users TABLE: %v", err)
		return
	}

	defer func() {
		if _, err := DB.Exec("DELETE FROM users WHERE id=$1", checkUserID); err != nil {
			t.Errorf("DELETE Error: %v. FIX THE DB MANUALY!", err)
			return
		}
	}()

	user := models.User{}

	if err := DB.Get(&user, "SELECT * FROM users WHERE id=$1", checkUserID); err != nil {
		t.Errorf("Error on db.Get: %v", err)
		return
	}

	if !userEqual(user, want) {
		t.Errorf("SELECT Error: db.Get() = %v, want = %v", user, want)
	}
}

func normalize(tm time.Time) (time.Time, error) {
	return time.Parse(
		"2006-01-02T15:04:05+07:00", tm.String(),
	)
}

func userEqual(created, want models.User) bool {
	created.RegistrationTime, _ = normalize(created.RegistrationTime)
	want.RegistrationTime, _ = normalize(want.RegistrationTime)

	return reflect.DeepEqual(created, want)
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

func TestGetUser(t *testing.T) {
	_, err := DB.NamedExec(insertUserRequest, testUser)
	if err != nil {
		t.Errorf("Error on TestGetUser: can't insert testUser manualy: %v", err)
		return
	}
	defer func() {
		DB.MustExec("DELETE FROM users WHERE id=$1", testUser.Id)
	}()

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			"GetUser() Test",
			args{checkUserID},
			&testUser,
			false,
		},
		{
			"GetUser() Test Error",
			args{uuid.New().String()},
			&testUser,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !userEqual(*got, *tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVideo(t *testing.T) {
	_, err := DB.NamedExec(
		"INSERT INTO videos VALUES (:id, :title, :localization, :upload_time, :file_path, :likes_count, :views_count)",
		testVideo,
	)
	if err != nil {
		t.Errorf("Error on TestGetVideo: can't insert testVideo manualy")
		return
	}
	defer func() {
		DB.MustExec("DELETE FROM videos WHERE id=$1", testVideo.Id)
	}()
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Video
		wantErr bool
	}{
		{
			"GetVideo() Test",
			args{checkVideoId},
			&testVideo,
			false,
		},
		{
			"GetVideo() Test Error",
			args{uuid.New().String()},
			&testVideo,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetVideo(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVideo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
