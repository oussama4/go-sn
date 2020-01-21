package models

import "time"

type User struct {
	ID        int
	Username  string `db:"user_name"`
	Email     string
	Password  string
	Bio       string
	Avatar    string
	CreatedAt time.Time `db:"created_at"`
}
