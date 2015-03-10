package controllers

import (
	"github.com/konek/auth-server/models"
	"github.com/konek/auth-server/tools"
)

// EditQuery ...
type EditQuery struct {
	Username  string                 `json:"username"`
	Password  string                 `json:"password"`
	Domains   []string               `json:"domains"`
	Enable    bool                   `json:"enable"`
	Variables map[string]interface{} `json:"variables"`
}

// EditResponse ...
type EditResponse struct {
	Status string `json:"status"`
}

// EditUser update user. Enable will be set to true if not specified. Password and Salt are updated if necessary.
func EditUser(handle tools.Handle) (interface{}, error) {
	var user models.User
	q := EditQuery{
		Enable: true,
	}

	uid := handle.P.ByName("uid")
	if uid == "" {
		return nil, tools.NewError(nil, 400, "bad request: missing userID")
	}
	if len(uid) != 24 {
		return nil, tools.NewError(nil, 400, "bad request: invalid userID")
	}

	err := tools.ParseBody(handle.R.Body, &q)
	if err != nil {
		return nil, tools.NewError(err, 400, "bad request: couldn't parse body")
	}

	user.IDFromHex(uid)
	if q.Username != "" {
		user.Username = q.Username
	}
	if q.Password != "" {
		if len(q.Password) < handle.C.PasswordMinLength {
			return nil, tools.NewError(nil, 400, "bad request: password is too short")
		}
		user.Password = q.Password
	}
	if q.Domains != nil {
		user.Domains = q.Domains
	}
	if q.Variables != nil {
		user.Variables = q.Variables
	}
	user.Enable = q.Enable

	err = user.Update()
	if err != nil {
		return nil, err
	}

	return EditResponse{
		Status: "ok",
	}, nil
}
