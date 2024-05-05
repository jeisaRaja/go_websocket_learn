package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID `db:"id" json:"id"`
	Username       string    `db:"username" json:"username"`
	Email          string    `db:"email" json:"email"`
	Password       password  `db:"password" json:"-"`
	ProfilePicture string    `db:"profile_picture" json:"profile_picture"`
	AccessToken    string    `db:"-" json:"token"`
}

type password struct {
	plainText *string
	hash      []byte
}

type AuthSignup struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthSignin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 10)
	if err != nil {
		return err
	}
	p.plainText = &plaintext
	p.hash = hash
	return nil
}

func (p *password) Match(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (p *password) GetHash() []byte {
	return p.hash
}
