package models

import (
	"log"

	"github.com/gocraft/dbr/v2"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type pgUserStore struct {
	logger *log.Logger
	db     *dbr.Connection
}

func NewUserStore(l *log.Logger, db *dbr.Connection) UserStore {
	return &pgUserStore{l, db}
}

func (ps *pgUserStore) Insert(name, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := User{
		Username: name,
		Email:    email,
		Password: string(hashed),
	}
	sess := ps.db.NewSession(nil)

	_, err = sess.InsertInto("users").Columns("user_name", "email", "password").Record(&user).Exec()
	if err != nil {
		// TODO: doesn't work as expected
		if pqError, ok := err.(pq.Error); ok {
			if pqError.Code == "23505" {
				ps.logger.Printf("column name: %s", pqError.Column)
				if pqError.Column == "user_name" {
					return ErrUsernameExist
				} else if pqError.Column == "email" {
					return ErrEmailExist
				}
			}
		}
		return err
	}
	return nil
}

func (ps *pgUserStore) Get(id int) (*User, error) {
	user := &User{}
	sess := ps.db.NewSession(nil)
	err := sess.Select("*").From("users").Where("id=?", id).LoadOne(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ps *pgUserStore) Authenticate(email, password string) (int, error) {
	user := &User{}
	sess := ps.db.NewSession(nil)
	err := sess.Select("*").From("users").Where("email=?", email).LoadOne(user)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return user.ID, nil
}
