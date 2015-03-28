package controllers

import (
	"time"

	"go.konek.io/auth-server/models"
	"go.konek.io/auth-server/tools"
	"go.konek.io/mgo"
)

// InfosSessionResponse ...
type InfosSessionResponse struct {
	Status  string          `json:"status"`
	Session SessionResponse `json:"session"`
}

// InfosSession get infos about the session (including remaining time)
func InfosSession(handle tools.Handle, db *mgo.DbQueue) (interface{}, error) {
	var resp InfosSessionResponse
	var session models.Session

	sid := handle.P.ByName("token")
	if sid == "" {
		return nil, tools.NewError(nil, 400, "bad request: missing token")
	}
	if tools.CheckID(sid) == false {
		return nil, tools.NewError(nil, 400, "bad request: invalid token")
	}

	session.IDFromHex(sid)
	err := session.Get(db)
	if err != nil {
		return nil, err
	}

	remaining := int(session.Expire - time.Now().Unix())
	if remaining <= 0 {
		return nil, tools.NewError(nil, 404, "not found: session is expired")
	}
	resp.Status = "ok"
	resp.Session.UserID = session.UserID.Hex()
	resp.Session.Domain = session.Domain
	resp.Session.Expire = session.Expire
	resp.Session.Remaining = remaining
	return resp, nil
}
