package models

type Tag struct {
	Id   uint   `db:"id"`
	Name string `db:"name"`
}
