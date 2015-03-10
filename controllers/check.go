package controllers

import (
	"time"

	"bitbucket.org/konek/auth-server/models"
	"bitbucket.org/konek/auth-server/tools"
)

// CheckRequest ...
type CheckRequest struct {
	Token  string `json:"token"`
	Domain string `json:"domain"`
}

// CheckResponse ...
type CheckResponse struct {
	Status  string          `json:"status"`
	Session SessionResponse `json:"session"`
}

// Check if a session is expired, and if it grants access to the specified domain
func Check(handle tools.Handle) (interface{}, error) {
	var q CheckRequest
	var resp CheckResponse
	var session models.Session

	err := tools.ParseBody(handle.R.Body, &q)
	if err != nil {
		return nil, tools.NewError(err, 400, "bad request: couldn't parse body")
	}

	if q.Token == "" {
		return nil, tools.NewError(err, 400, "bad request: token is missing")
	}
	if tools.CheckID(q.Token) == false {
		return nil, tools.NewError(err, 400, "bad request: invalid token")
	}
	if q.Domain == "" {
		return nil, tools.NewError(err, 400, "bad request: domain is missing")
	}
	if q.Domain == "/" {
		return nil, tools.NewError(nil, 400, "bad request: illegal domain")
	}

	session.IDFromHex(q.Token)
	err = session.Get()
	if err != nil {
		return nil, err
	}

	if session.Expire < time.Now().Unix() {
		return nil, tools.NewError(nil, 404, "not found: session is expired")
	}

	if tools.CheckDomain(q.Domain, session.Domain) == false {
		return nil, tools.NewError(nil, 403, "forbidden: restricted domain")
	}

	resp.Status = "ok"
	resp.Session.UserID = session.UserID.Hex()
	resp.Session.Expire = session.Expire
	resp.Session.Remaining = int(session.Expire - time.Now().Unix())

	return resp, nil
}
