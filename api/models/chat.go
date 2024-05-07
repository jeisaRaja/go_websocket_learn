package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID      uuid.UUID `json:"id" db:"id"`
	Message string    `json:"message" db:"message"`
	Room    string    `json:"room" db:"room"`
	From    string    `json:"from" db:"from_user"`
	Sent    time.Time `json:"sent" db:"sent"`
}
