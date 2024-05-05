package database

import (
	"database/sql"
	"jeisaraja/websocket_learn/models"

)

type Queries struct {
	DB *sql.DB
}

func (q *Queries) CreateUser(u *models.User) error {
	query := "INSERT INTO users (id, username, email, password, profile_picture) VALUES ($1,$2,$3,$4,$5)"
	_, err := q.DB.Exec(query, u.ID, u.Username, u.Email, u.Password.GetHash(), u.ProfilePicture)
	if err != nil {
		return err
	}
	return nil
}

func (q *Queries) GetUserPassword(username string) *[]byte {
	query := "SELECT password FROM users WHERE id = $1"
	row := q.DB.QueryRow(query, username)

	var password []byte
	err := row.Scan(&password)
	if err == sql.ErrNoRows {
		return nil
	}
	return &password
}
