package models

import (
	"time"

	"github.com/google/uuid"
)

type JwtToken struct {
	Uuid         uuid.UUID
	Token        string
	CreationTime time.Time
}
