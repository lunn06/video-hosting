package database

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lunn06/video-hosting/internal/config"
	"github.com/lunn06/video-hosting/internal/models"
)

func mustTestInit() {
	config.CFG = config.MustLoad("../../configs/main.yaml")
	DB = MustCreate(config.CFG)
}

var (
	checkUserId      uint32 = 0
	checkEmail              = "check@tt.com"
	checkChannelName        = "test"
	checkPassword           = "1234"
	checkTime               = time.Now()

	testUser = models.User{
		Id:               checkUserId,
		Email:            checkEmail,
		ChannelName:      checkChannelName,
		Password:         checkPassword,
		RegistrationTime: checkTime,
	}

	checkVideoUuid         = uuid.New()
	checkTitle             = "Test Title"
	checkLocalization      = "ru_RU"
	checkUploadTime        = time.Now()
	checkFilePath          = "/some/test/dir/to/file.mp4"
	checkLikesCount   uint = 100
	checkViewsCount   uint = 10000

	testVideo = models.Video{
		Uuid:         checkVideoUuid,
		Title:        checkTitle,
		Localization: checkLocalization,
		UploadTime:   checkUploadTime,
		FilePath:     checkFilePath,
		LikesCount:   checkLikesCount,
		ViewsCount:   checkViewsCount,
	}

	testInsertUserRequest = `INSERT INTO users VALUES (
		:id, :email, :channel_name, :password, :registration_time
		)`

	testInsertVideoRequest = `INSERT INTO videos VALUES (
		:uuid, :title, :localization, :upload_time, :file_path, :likes_count, :views_count
		)`
)

func TestMustCreate(t *testing.T) {
	mustTestInit()
	want := testUser

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Error on database.MustCreate(): %v", r)
			return
		}
	}()

	t.Run("MustCreate() Test", func(t *testing.T) {
		tx := DB.MustBegin()
		if _, err := tx.Exec(
			"INSERT INTO users VALUES ($1, $2, $3, $4, $5)",
			checkUserId, checkEmail, checkChannelName, checkPassword, checkTime,
		); err != nil {
			t.Errorf("INSERT test values in users TABLE error = %v", err)
			return
		}

		if err := tx.Commit(); err != nil {
			t.Errorf("Error on COMMIT in users TABLE: %v", err)
			return
		}

		defer func() {
			if _, err := DB.Exec("DELETE FROM users WHERE email=$1", checkEmail); err != nil {
				t.Errorf("DELETE Error: %v. FIX THE DB MANUALY!", err)
				return
			}
		}()

		user := models.User{}

		if err := DB.Get(&user, "SELECT * FROM users WHERE email=$1", checkEmail); err != nil {
			t.Errorf("Error on db.Get: %v", err)
			return
		}

		if !userEqual(user, want) {
			t.Errorf("SELECT Error: db.Get() = %v, want = %v", user, want)
		}
	})
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

func videoEqual(created, want models.Video) bool {
	created.UploadTime, _ = normalize(created.UploadTime)
	want.UploadTime, _ = normalize(want.UploadTime)

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
			"getPgAddress() Test",
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
	mustTestInit()

	randomUUID := rand.Uint32()
	_, err := DB.NamedExec(testInsertUserRequest, testUser)
	if err != nil {
		t.Errorf("TestGetUser() error: can't insert testUser manualy = %v", err)
		return
	}

	defer func() {
		DB.MustExec("DELETE FROM users WHERE id=$1", testUser.Id)
	}()

	type args struct {
		id uint32
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			"GetUser() Test",
			args{checkUserId},
			&testUser,
			false,
		},
		{
			"GetUser() Test Error",
			args{randomUUID},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.id)

			//DB.MustExec("DELETE FROM users WHERE id=$1", testUser.Uuid)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil && got != nil {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
				return
			}
			if tt.want == nil && got == nil {
				return
			}
			if !userEqual(*got, *tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
			if r := recover(); r != nil {
				DB.MustExec("DELETE FROM users WHERE id=$1", testUser.Id)
			}
		})
	}
}

func TestGetVideo(t *testing.T) {
	mustTestInit()

	_, err := DB.NamedExec(testInsertVideoRequest, testVideo)
	if err != nil {
		t.Errorf("Error on TestGetVideo: can't insert testVideo manualy")
		return
	}
	defer func() {
		DB.MustExec("DELETE FROM videos WHERE uuid=$1", testVideo.Uuid)
	}()
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Video
		wantErr bool
	}{
		{
			"GetVideo() Test",
			args{checkVideoUuid},
			&testVideo,
			false,
		},
		{
			"GetVideo() Test Error",
			args{uuid.New()},
			nil,
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
			if tt.want == nil && got != nil {
				t.Errorf("GetVideo() got = %v, want %v", got, tt.want)
				return
			}
			if tt.want == nil && got == nil {
				return
			}
			if !videoEqual(*got, *tt.want) {
				t.Errorf("GetVideo() got = %v, want %v", got, tt.want)
			}
			if r := recover(); r != nil {
				DB.MustExec("DELETE FROM videos WHERE uuid=$1", testVideo.Uuid)
			}
		})
	}
}

func TestInsertUser(t *testing.T) {
	mustTestInit()

	insertUser := models.User{
		Id:    uint32(rand.Int31()),
		Email: "insert@mail.ru",
	}

	_, err := DB.NamedExec(testInsertUserRequest, testUser)
	if err != nil {
		t.Errorf("Error on TestInsertUser: can't insert testUser manualy: %v", err)
		return
	}

	defer func() {
		DB.MustExec("DELETE FROM users WHERE email=$1 OR email=$2", testUser.Email, insertUser.Email)
	}()

	type args struct {
		user models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"InsertUser() Test",
			args{insertUser},
			false,
		},
		{
			"InsertUser() Test Error",
			args{testUser},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if r := recover(); r != nil {
				DB.MustExec("DELETE FROM users WHERE email=$1 OR email=$2", testUser.Email, insertUser.Email)
			}
		})
	}
}

func TestInsertVideo(t *testing.T) {
	mustTestInit()

	insertVideo := models.Video{
		Uuid:  uuid.New(),
		Title: "Insert Title",
	}

	_, err := DB.NamedExec(testInsertVideoRequest, testVideo)
	if err != nil {
		t.Errorf("Error on TestInsertVideo: can't insert testVideo manualy: %v", err)
		return
	}

	defer func() {
		DB.MustExec("DELETE FROM videos WHERE title=$1 OR title=$2", testVideo.Title, insertVideo.Title)
	}()

	type args struct {
		video models.Video
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"InsertVideo() Test",
			args{insertVideo},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertVideo(tt.args.video); (err != nil) != tt.wantErr {
				t.Errorf("InsertVideo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if r := recover(); r != nil {
				DB.MustExec("DELETE FROM videos WHERE title=$1 OR title=$2", testVideo.Title, insertVideo.Title)
			}
		})
	}
}
