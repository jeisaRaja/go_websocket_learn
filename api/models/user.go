package models

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID `db:"id" json:"id"`
	Username       string    `db:"username" json:"username"`
	Email          string    `db:"email" json:"email"`
	Password       string    `db:"password" json:"password"`
	ProfilePicture string    `db:"profile_picture" json:"profile_picture"`
	AccessToken    string    `json:"token"`
}
