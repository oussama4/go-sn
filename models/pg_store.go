package models

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"
)

var (
	ErrEmailExist = errors.New("Email allready exists")

	ErrUsernameExist = errors.New("username already exists")

	ErrInvalidCredentials = errors.New("Invalid credentials")
)

// New creates a new database connection
func New(l *log.Logger) (*dbr.Connection, error) {
	con, err := dbr.Open("postgres", os.Getenv("DATABASE_URL"), nil)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the databse")
	}
	err = con.Ping()
	if err != nil {
		return nil, err
	}

	l.Println("connected to database")
	return con, nil
}
