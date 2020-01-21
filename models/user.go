package models

import (
	"github.com/gocraft/dbr/v2"
	"time"
)

type User struct {
	ID        int
	Username  string `db:"user_name"`
	Email     string
	Password  string
	Bio       dbr.NullString
	Avatar    string
	CreatedAt time.Time `db:"created_at"`
}
