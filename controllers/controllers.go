package controllers

import (
	"bitbucket.org/konek/auth-server/tools"
)

// ControllerFunc ...
type ControllerFunc func(handle tools.Handle) (resp interface{}, err error)
