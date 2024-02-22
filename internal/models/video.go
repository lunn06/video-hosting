package models

import "time"

type Video struct {
	Id           string    `db:"id"`
	Title        string    `db:"title"`
	Localization string    `db:"localization"`
	UploadDate   time.Time `db:"upload_time"`
	FilePath     string    `db:"file_path"`
	LikesCount   uint      `db:"likes_count"`
	ViewsCount   uint      `db:"views_count"`
}
