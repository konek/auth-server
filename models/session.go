package models

import (
	"time"

	mgov2 "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"go.konek.io/auth-server/tools"
	"go.konek.io/mgo"
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
func (s *Session) Create(db *mgo.DbQueue, lifespan int64) (int, error) {
	s.ID = bson.NewObjectId()
	if s.Expire == 0 {
		s.Expire = time.Now().Unix() + lifespan
	}
	err := db.Push(func (db *mgo.Database, ec chan error) {
		ec <- db.C("sessions").Insert(s)
	})
	if err != nil {
		return 0, err
	}

	return int(lifespan), nil
}

// Delete a session in database
func (s *Session) Delete(db *mgo.DbQueue) error {
	var n int

	err := db.Push(func (db *mgo.Database, ec chan error) {
		var e error

		q := db.C("sessions").FindId(s.ID)
		n, e = q.Count()
		ec <- e
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return tools.NewError(nil, 404, "not found: session does not exist")
	}
	err = db.Push(func (db *mgo.Database, ec chan error) {
		ec <- db.C("sessions").RemoveId(s.ID)
	})
	return err
}

// Get a session from database
func (s *Session) Get(db *mgo.DbQueue) error {
	var n int

	err := db.Push(func (db *mgo.Database, ec chan error) {
		var e error

		q := db.C("sessions").FindId(s.ID)
		n, e = q.Count()
		if e != nil || n == 0 {
			ec <- e
			return
		}
		ec <- q.One(s)
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return tools.NewError(nil, 404, "not found: session does not exist")
	}
	return nil
}

// CleanSessions remove expired sessions older than `age`
func CleanSessions(db *mgo.DbQueue, age int64) (int, error) {
	var change *mgov2.ChangeInfo
	limit := time.Now().Unix() - age

	err := db.Push(func (db *mgo.Database, ec chan error) {
		var e error
		change, e = db.C("sessions").RemoveAll(bson.M{"expire": bson.M{"$lt": limit}})
		ec <- e
	})
	if err != nil {
		return 0, err
	}

	return change.Removed, nil
}

// ListSessions returns a list of all the sessions
func ListSessions(db *mgo.DbQueue) ([]Session, error) {
	var list []Session
	var n int

	err := db.Push(func (db *mgo.Database, ec chan error) {
		var e error

		q := db.C("sessions").Find(nil)
		n, e = q.Count()
		if e != nil || n == 0 {
			ec <- e
			return
		}
		list = make([]Session, n)
		ec <- q.All(&list)
	})

	if err != nil {
		return nil, err
	}
	return list, nil
}
