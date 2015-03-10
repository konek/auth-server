package controllers

import (
	"go.konek.io/auth-server/models"
	"go.konek.io/auth-server/tools"
)

// LogoffResponse ...
type LogoffResponse struct {
	Status string `json:"status"`
}

// Logoff deletes a session (expired or not)
func Logoff(handle tools.Handle) (interface{}, error) {
	var session models.Session

	sid := handle.P.ByName("token")
	if sid == "" {
		return nil, tools.NewError(nil, 400, "bad request: missing token")
	}
	if tools.CheckID(sid) == false {
		return nil, tools.NewError(nil, 400, "bad request: invalid token")
	}
	session.IDFromHex(sid)

	err := session.Delete()
	if err != nil {
		return nil, err
	}

	return LogoffResponse{"ok"}, nil
}
