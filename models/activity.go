package models

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/gocraft/dbr/v2"
)

// these constances define what types of activities can be created
const (
	CREATE = "create"
	LIKE   = "like"
	SHARE  = "share"
	OFFER  = "offer"
)

// these constances define what types of objects can be created
const (
	POST       = "post" // post contains both image and text
	TARGET     = "target"
	CONNECTION = "connection"
)

// ActivityRecord holds activity data data as it is stored in postgres
type ActivityRecord struct {
	Username string `db:"user_name"`
	Avatar   string
	AID      int    `db:"id"`
	Atype    string `db:"type"`
	Actor    int
	Content  []byte
	Created  time.Time
}

// BaseObject holds fields that are shared between all  object types
type BaseObject struct {
	ID    int
	OType string
}

type PostObject struct {
	BaseObject
	Txt string
	Img string
}

// BaseActivity holds data that is shared between all type of activities
type BaseActivity struct {
	ID      int
	AType   string
	Actor   int
	Created time.Time
}

type CreateActivity struct {
	BaseActivity
	Object PostObject
}

func NewCreateActivity(actor int, atype, txt, img string) *CreateActivity {
	return &CreateActivity{
		BaseActivity: BaseActivity{
			AType: atype,
			Actor: actor,
		},
		Object: PostObject{
			BaseObject: BaseObject{
				OType: POST,
			},
			Txt: txt,
			Img: img,
		},
	}
}

func NewCreateActivityFromRecord(a ActivityRecord) *CreateActivity {
	return &CreateActivity{
		BaseActivity: BaseActivity{
			ID:      a.AID,
			AType:   a.Atype,
			Actor:   a.Actor,
			Created: a.Created,
		},
	}
}

func (ca *CreateActivity) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&ca.Object); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (ca *CreateActivity) Decode(b []byte) error {
	r := bytes.NewReader(b)
	if err := gob.NewDecoder(r).Decode(&ca.Object); err != nil {
		return err
	}
	return nil
}

func (ca *CreateActivity) Save(db *dbr.Connection) error {
	c, err := ca.Encode()
	if err != nil {
		return err
	}

	sess := db.NewSession(nil)
	_, err = sess.InsertInto("activities").
		Pair("type", ca.AType).
		Pair("actor", ca.Actor).
		Pair("content", c).Exec()
	if err != nil {
		return err
	}

	return nil
}

// create a slice of Activity types from ActivityRecord types for the purpuse of encoding it to json
func activitiesFromRecords(activities []ActivityRecord) ([]interface{}, error) {
	r := []interface{}{}

	for _, v := range activities {
		switch v.Atype {
		case CREATE:
			a := NewCreateActivityFromRecord(v)
			if err := a.Decode(v.Content); err != nil {
				return nil, err
			}
			r = append(r, a)
		}
	}
	return r, nil
}
