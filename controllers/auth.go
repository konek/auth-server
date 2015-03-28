package controllers

import (
	"go.konek.io/auth-server/models"
	"go.konek.io/auth-server/tools"
	"go.konek.io/rest"
	"go.konek.io/mgo"
)

// AuthRequest ...
type AuthRequest struct {
	Domain   string `json:"domain"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse ...
type AuthResponse struct {
	Status string `json:"status"`
	UserID string `json:"uid"`
}

// Auth authenticate a user, no session is created
func Auth(handle tools.Handle, db *mgo.DbQueue) (interface{}, error) {
	var q AuthRequest
	var user models.User

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
	return AuthResponse{
		Status: "ok",
		UserID: user.ID.Hex(),
	}, nil
}
