package controllers

import (
	"bitbucket.org/konek/mgo"
	"go.konek.io/auth-server/tools"
)

// ControllerFunc ...
type ControllerFunc func(handle tools.Handle, db *mgo.DbQueue) (resp interface{}, err error)
