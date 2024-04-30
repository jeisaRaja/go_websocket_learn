package models

import "time"

type Chat struct {
	Message string    `json:"message"`
	From    string    `json:"from"`
	Room    string    `json:"room"`
	Sent    time.Time `json:"sent"`
}
