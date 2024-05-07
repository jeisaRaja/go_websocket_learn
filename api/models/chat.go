package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Message  string    `json:"message" db:"message"`
	Room     string    `json:"room" db:"room"`
	FromID   uuid.UUID `json:"from_id" db:"from_user_id"`
	FromName string    `json:"from_name" db:"from_user_name"`
	Sent     time.Time `json:"sent" db:"sent"`
}
