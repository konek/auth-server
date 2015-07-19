package controllers

import (
	"gopkg.in/konek/mgo.v1"
	"gopkg.in/konek/auth-server.v1/models"
	"gopkg.in/konek/auth-server.v1/tools"
	"go.konek.io/rest"

	"github.com/asaskevich/govalidator"
)

// LoginRequest ...
type LoginRequest struct {
	Domain   string `json:"domain"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse ...
type LoginResponse struct {
	Status  string          `json:"status"`
	Session SessionResponse `json:"session"`
}

// Login a user, creating a new session.
func Login(handle tools.Handle, db *mgo.DbQueue) (interface{}, error) {
	var q LoginRequest
	var user models.User
	var session models.Session
	var resp LoginResponse

	err := rest.Parse(handle.R, &q)
	if err != nil {
		return nil, tools.NewError(err, 400, "bad request: couldn't parse body")
	}

	if q.Domain == "" {
		return nil, tools.NewError(nil, 400, "bad request: domain is missing")
	}
	if q.Domain == "/" {
		return nil, tools.NewError(nil, 400, "bad request: illegal domain")
	}

	if q.Username == "" {
		return nil, tools.NewError(nil, 400, "bad request: username is missing")
	}
	if q.Password == "" {
		return nil, tools.NewError(nil, 400, "bad request: password is missing")
	}

	user.Username = q.Username
	user.Password = q.Password

	if govalidator.IsEmail(user.Username) == false {
		return nil, tools.NewError(nil, 400, "bad request: username must be a valid email")
	}

	user.Username, err = govalidator.NormalizeEmail(user.Username)
	if err != nil {
		return nil, tools.NewError(nil, 400, "bad request: username must be a valid email")
	}

	ok, err := user.Check(db)
	if err != nil {
		return nil, err
	}
	if ok == false {
		return nil, tools.NewError(nil, 403, "forbidden: invalid user or password")
	}
	if user.Enable == false {
		return nil, tools.NewError(nil, 403, "forbidden: user is diabled")
	}
	ok = user.CheckDomain(q.Domain)
	if ok == false {
		return nil, tools.NewError(nil, 403, "forbidden: restricted domain")
	}

	session.UserID = user.ID
	session.Domain = q.Domain
	remaining, err := session.Create(db, handle.C.SessionLifespan)
	if err != nil {
		return nil, err
	}
	resp.Status = "ok"
	resp.Session.Token = session.ID.Hex()
	resp.Session.UserID = session.UserID.Hex()
	resp.Session.Expire = session.Expire
	resp.Session.Remaining = remaining
	return resp, nil
}
