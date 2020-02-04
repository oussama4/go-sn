package models

type UserStore interface {
	Insert(name, email, password string) error
	Get(id int) (*User, error)
	Authenticate(email, password string) (int, error)
	Update(id int, updated map[string]string) error
}

type ConnectionStore interface {
	AreConnected(is, with int) (bool, error)
	// get the IDs of  a user connections
	Connections(user int) ([]int, error)
}

type ActivityStore interface {
	// given  a slice of user IDs  return a slice of activities
	Activities(users []int, offset, limit int) ([]interface{}, error)
}
