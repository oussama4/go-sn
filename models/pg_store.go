package models

import (
	"fmt"
	"log"
	"os"

	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// New creates a new database connection
func New(l *log.Logger) (*dbr.Connection, error) {
	con, err := dbr.Open("postgres", os.Getenv("DATABASE_URL"), nil)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the databse")
	}

	l.Println("connected to database")
	return con, nil
}

type pgUserStore struct {
	logger *log.Logger
	db     *dbr.Connection
}

func NewUserStore(l * log.Logger, db *dbr.Connection) UserStore  {
	return &pgUserStore{l, db}
}

func (ps *pgUserStore) Insert(name, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user := User{
		Username: name,
		Email: email,
		Password: string(hashed),
	}
	sess := ps.db.NewSession(nil)
	_, err := sess.InsertInto("users").Columns("user_name", "email", "password").Record(&user).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (ps *pgUserStore) Get(id int) (*User, error)  {
	user := &User{}
	sess := ps.db.NewSession(nil)
	err := sess.Select("*").From("users").Where("id=$1", id).LoadOne(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
