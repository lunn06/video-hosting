package models

import (
	"time"

	"github.com/google/uuid"
)

type Video struct {
	Uuid         uuid.UUID `db:"uuid"`
	Title        string    `db:"title"`
	Localization string    `db:"localization"`
	UploadTime   time.Time `db:"upload_time"`
	FilePath     string    `db:"file_path"`
	LikesCount   uint      `db:"likes_count"`
	ViewsCount   uint      `db:"views_count"`

	Tags  []*Tag
	Owner *User
}
