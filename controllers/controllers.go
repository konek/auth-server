package controllers

import (
	"go.konek.io/auth-server/tools"
	"go.konek.io/mgo"
)

// ControllerFunc ...
type ControllerFunc func(handle tools.Handle, db *mgo.DbQueue) (resp interface{}, err error)
