package models

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/konek/auth-server/tools"
)

// Session is the model for sessions
type Session struct {
	ID     bson.ObjectId `bson:"_id" json:"sid"`
	UserID bson.ObjectId `json:"uid"`
	Domain string        `json:"domain"`
	Expire int64         `json:"expire"`
}

// IDFromHex fill the ID from a string hex
func (s *Session) IDFromHex(hex string) {
	s.ID = bson.ObjectIdHex(hex)
}

// Create a new session in database
func (s *Session) Create(lifespan int64) (int, error) {
	s.ID = bson.NewObjectId()
	if s.Expire == 0 {
		s.Expire = time.Now().Unix() + lifespan
	}
	err := gDb.C("sessions").Insert(s)
	if err != nil {
		return 0, err
	}

	return int(lifespan), nil
}

// Delete a session in database
func (s *Session) Delete() error {
	q := gDb.C("sessions").FindId(s.ID)
	n, err := q.Count()
	if err != nil {
		return err
	}
	if n == 0 {
		return tools.NewError(nil, 404, "not found: session does not exist")
	}
	err = gDb.C("sessions").RemoveId(s.ID)
	return err
}

// Get a session from database
func (s *Session) Get() error {
	q := gDb.C("sessions").FindId(s.ID)
	n, err := q.Count()
	if err != nil {
		return err
	}
	if n == 0 {
		return tools.NewError(nil, 404, "not found: session does not exist")
	}
	err = q.One(s)
	if err != nil {
		return err
	}
	return nil
}

// CleanSessions remove expired sessions older than `age`
func CleanSessions(age int64) (int, error) {
	limit := time.Now().Unix() - age

	fmt.Println(age)
	fmt.Println(limit)
	change, err := gDb.C("sessions").RemoveAll(bson.M{"expire": bson.M{"$lt": limit}})
	if err != nil {
		return 0, err
	}

	return change.Removed, nil
}

// ListSessions returns a list of all the sessions
func ListSessions() ([]Session, error) {
	var list []Session
	q := gDb.C("sessions").Find(nil)

	n, err := q.Count()
	if err != nil {
		return nil, err
	}
	list = make([]Session, n)
	err = q.All(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
