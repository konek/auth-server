package controllers

import (
	"bitbucket.org/konek/auth-server/models"
	"bitbucket.org/konek/auth-server/tools"
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
func Auth(handle tools.Handle) (interface{}, error) {
	var q LoginRequest
	var user models.User

	err := tools.ParseBody(handle.R.Body, &q)
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

	ok, err := user.Check()
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
