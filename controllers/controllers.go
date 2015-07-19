package controllers

import (
	"bitbucket.org/konek/mgo"
	"gopkg.in/konek/auth-server.v1/tools"
)

// ControllerFunc ...
type ControllerFunc func(handle tools.Handle, db *mgo.DbQueue) (resp interface{}, err error)
