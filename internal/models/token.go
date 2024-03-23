package models

import (
	"time"

	"github.com/google/uuid"
)

type JwtToken struct {
	Uuid         uuid.UUID `db:"uuid"`
	Token        string    `db:"token"`
	CreationTime time.Time `db:"creation_time"`
}
