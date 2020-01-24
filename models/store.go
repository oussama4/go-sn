package models

type UserStore interface {
	Insert(name, email, password string) error
	Get(id int) (*User, error)
	Authenticate(email, password string) (int, error)
}

type ConnectionStore interface {
	AreConnected(is, with int) (bool, error)
}
