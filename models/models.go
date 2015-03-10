package models

import (
	"log"

	"gopkg.in/mgo.v2"

	"bitbucket.org/konek/auth-server/config"
)

var gSession *mgo.Session
var gDb *mgo.Database

// Init dial the connection to the mongodb instance
func Init(conf config.Conf) error {
	gSession, err := mgo.Dial(conf.DbURL)
	if err != nil {
		log.Println("couldn't connect to database :", err, ", mongo_url is", conf.DbURL)
		return err
	}
	gDb = gSession.DB(conf.DbName)
	return nil
}
