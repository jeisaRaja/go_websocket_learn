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
