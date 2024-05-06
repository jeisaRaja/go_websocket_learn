package database

import (
	"database/sql"
	"fmt"
	"jeisaraja/websocket_learn/models"

	"github.com/google/uuid"
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
	query := "SELECT password FROM users WHERE username = $1"
	row := q.DB.QueryRow(query, username)

	var password []byte
	err := row.Scan(&password)
	if err == sql.ErrNoRows {
		fmt.Println("No row found")
		return nil
	}
	return &password
}

func (q *Queries) InsertChat(chat *models.Chat) error {
	query := "INSERT INTO chats (id, message, room, from_user,sent) VALUES ($1,$2,$3,$4,$5)"
	fromID, err := q.GetUserID(chat.From)
	if err != nil {
		return err
	}
	id := uuid.New()
	_, err = q.DB.Exec(query, id, chat.Message, chat.Room, fromID, chat.Sent)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(query)
		return err
	}
	return nil
}

func (q *Queries) GetUserID(username string) (string, error) {
	query := "SELECT id FROM users WHERE username = $1"
	var uname string
	result := q.DB.QueryRow(query, username)
	err := result.Scan(&uname)
	if err != nil {
		return "", err
	}
	return uname, nil
}
