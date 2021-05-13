package models

type Connection struct {
	UserOne int `db:"user_one"`
	UserTwo int `db:"user_two"`
	Status  int
}
