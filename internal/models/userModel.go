package models

type User struct {
	Id          string `db:"id"`
	Email       string `db:"email"`
	ChannelName string `db:"channel_name"`
	Password    string `db:"password"`
}
