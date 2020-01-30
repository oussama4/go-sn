package models

import (
	"bytes"
	"encoding/gob"

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
	ID    int
	AType string
	Actor int
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

func (ca *CreateActivity) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(ca.Object); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (ca *CreateActivity) Save(db *dbr.Connection) error {
	c, err := ca.Serialize()
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
