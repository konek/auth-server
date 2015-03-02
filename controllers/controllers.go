
package controllers

import (
  "github.com/konek/auth-server/tools"
)

// ControllerFunc ...
type ControllerFunc func (handle tools.Handle) (resp interface{}, err error)

