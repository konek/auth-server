package controllers

import (
	"gopkg.in/konek/mgo.v1"
	"gopkg.in/konek/auth-server.v1/tools"
)

// ControllerFunc ...
type ControllerFunc func(handle tools.Handle, db *mgo.DbQueue) (resp interface{}, err error)
