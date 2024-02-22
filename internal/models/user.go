package models

import (
	"time"
)

type User struct {
	Id               string    `db:"id"`
	Email            string    `db:"email"`
	ChannelName      string    `db:"channel_name"`
	Password         string    `db:"password"`
	RegistrationTime time.Time `db:"data_reg"`
}
