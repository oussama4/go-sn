package models

import "time"

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	Bio       string
	Avatar    string
	CreatedAt time.Time
}
