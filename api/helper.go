package main

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func ComparePasswords(hashedPassword []byte, plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(plaintextPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
