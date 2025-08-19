package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	// getUserByID is a SQL string to select a user by their ID. It returns the user's ID, created_at, username, and hashed_password.
	getUserByID string = `SELECT id, created_at, username, hashed_password FROM users WHERE id = $1;`
	// getUserByUsername is a SQL string to select a user by their username. It returns the user's ID, created_at, username, and hashed_password.
	getUserByUsername string = `SELECT id, created_at, username, hashed_password FROM users WHERE username = $1;`
	// saveUser is a SQL string to insert into the users table. It requires the username and hashed_password as input and returns the id of the newly created user.
	saveUser string = `INSERT INTO users (username, hashed_password) VALUES ($1, $2) RETURNING id;`
)

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	Username       string
	HashedPassword string
}

func (db *DB) GetUserByID(id uuid.UUID) (*User, error) {
	row := db.conn.QueryRow(context.Background(), getUserByID, id)
	var user User
	err := row.Scan(&user.ID, &user.CreatedAt, &user.Username, &user.HashedPassword)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByUsername(username string) (*User, error) {
	row := db.conn.QueryRow(context.Background(), getUserByUsername, username)
	var user User
	err := row.Scan(&user.ID, &user.CreatedAt, &user.Username, &user.HashedPassword)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) InsertUser(username, hashedPassword string) (id uuid.UUID, err error) {
	row := db.conn.QueryRow(context.Background(), saveUser, username, hashedPassword)
	err = row.Scan(&id)
	return
}
