package models

import "time"

type Chat struct {
  Message string    `json:"message" db:"message"`
  Room    string    `json:"room" db:"room"`
  From    string    `json:"from" db:"from"`
  Sent    time.Time `json:"sent" db:"sent"`
}
