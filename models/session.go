package models

import (
	"time"

	mgov2 "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/konek/mgo.v1"
	"gopkg.in/konek/auth-server.v1/tools"
)

// Session is the model for sessions
type Session struct {
	ID     bson.ObjectId `bson:"_id" json:"_id"`
	UserID bson.ObjectId
	Domain string
	Expire int64
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
	err := db.Insert("sessions", s)
	if err != nil {
		return 0, err
	}
	return int(lifespan), nil
}

// Delete a session in database
func (s *Session) Delete(db *mgo.DbQueue) error {
	n, err := db.Count("sessions", mgo.M{"_id": s.ID})
	if err != nil {
		return err
	}
	if n == 0 {
		return tools.NewError(nil, 404, "not found: session does not exist")
	}
	err = db.RemoveID("sessions", s.ID)
	return err
}

// Get a session from database
func (s *Session) Get(db *mgo.DbQueue) error {
	err := db.FindOneID("sessions", s, s.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s Session) Check(domain string) bool {
	if s.Expire < time.Now().Unix() {
		return false
	}
	if tools.CheckDomain(domain, s.Domain) == false {
		return false
	}
	return true
}

// CleanSessions remove expired sessions older than `age`
func CleanSessions(db *mgo.DbQueue, age int64) (int, error) {
	var change *mgov2.ChangeInfo
	limit := time.Now().Unix() - age

	err := db.Push(func(db *mgo.Database, ec chan error) {
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
	err := db.Find("sessions", &list, nil)
	if err != nil {
		return nil, err
	}
	return list, nil
}
