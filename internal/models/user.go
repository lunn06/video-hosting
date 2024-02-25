package models

import (
	"time"
)

type User struct {
	Id               uint32    `db:"id"`
	Email            string    `db:"email"`
	ChannelName      string    `db:"channel_name"`
	Password         string    `db:"password"`
	RegistrationTime time.Time `db:"registration_time"`

	Videos []*Video
	Role   *Role
}
