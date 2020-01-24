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
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}
