package models

type Tag struct {
	Id   uint32 `db:"id"`
	Name string `db:"name"`
}
