package models

type UserStore interface {
	Insert(name, email, password string) error
	Get(id int) (*User, error)
}
