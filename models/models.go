package models

import (
	"gopkg.in/mgo.v2"

	"go.konek.io/auth-server/config"
)

var gSession *mgo.Session
var gDb *mgo.Database

// Init dial the connection to the mongodb instance
func Init(conf config.Conf) error {
	gSession, err := mgo.Dial(conf.DbURL)
	if err != nil {
		return err
	}
	gDb = gSession.DB(conf.DbName)
	return nil
}
