package models

import (
	"log"

	"github.com/gocraft/dbr/v2"
)

type pgActivityStore struct {
	logger *log.Logger
	db     *dbr.Connection
}

func NewActivityStore(l *log.Logger, db *dbr.Connection) ActivityStore {
	return &pgActivityStore{l, db}
}

func (as *pgActivityStore) Activities(ids []int, offset, limit int) ([]interface{}, error) {
	sess := as.db.NewSession(nil)
	r := []ActivityRecord{}

	_, err := sess.Select("users.user_name", "users.avatar", "activities.*").
		From("activities").Join("users", "activities.actor=users.id").
		Where("activities.actor in ?", ids).OrderDesc("created").
		Offset(uint64(offset)).Limit(uint64(limit)).Load(&r)
	if err != nil {
		return nil, err
	}
	acts, err := activitiesFromRecords(r)
	if err != nil {
		return nil, err
	}
	return acts, nil
}
