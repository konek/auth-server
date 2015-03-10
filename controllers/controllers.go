package controllers

import (
	"go.konek.io/auth-server/tools"
)

// ControllerFunc ...
type ControllerFunc func(handle tools.Handle) (resp interface{}, err error)
