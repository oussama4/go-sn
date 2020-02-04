package models

import (
	"log"

	"github.com/gocraft/dbr/v2"
)

type pgConnStore struct {
	logger *log.Logger
	db     *dbr.Connection
}

func NewConnStore(l *log.Logger, db *dbr.Connection) ConnectionStore {
	return &pgConnStore{l, db}
}

// checks if two users are connected
func (cs *pgConnStore) AreConnected(is, with int) (bool, error) {
	q := `select * from connections where user_one=? and user_two=? and status=?
	union select * from connections where user_one=? and user_two=? and status=?`
	var conns []Connection
	sess := cs.db.NewSession(nil)
	c, err := sess.SelectBySql(q, is, with, 1, with, is, 1).Load(&conns)
	if err != nil {
		return false, err
	} else if c == 0 {
		return false, nil
	}
	return true, nil
}

func (cs *pgConnStore) Connections(user int) ([]int, error) {
	q := `select user_one as u from connections where user_two=? and status=? 
	union select user_two as u from connections where user_one=? and status=?`
	r := []int{}
	sess := cs.db.NewSession(nil)
	_, err := sess.SelectBySql(q, user, 1, user, 1).Load(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}
