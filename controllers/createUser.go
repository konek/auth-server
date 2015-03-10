package controllers

import (
	"bitbucket.org/konek/auth-server/models"
	"bitbucket.org/konek/auth-server/tools"
)

// CreateResponse ...
type CreateResponse struct {
	Status string `json:"status"`
	UserID string `json:"uid"`
}

// CreateUser create a new user. Checks for duplicate users and password-length requirement
func CreateUser(handle tools.Handle) (interface{}, error) {
	var user models.User

	user.Enable = true
	user.Domains = nil
	user.Variables = nil
	err := tools.ParseBody(handle.R.Body, &user)
	if err != nil {
		return nil, tools.NewError(err, 400, "bad request: couldn't parse body")
	}

	if user.Username == "" {
		return nil, tools.NewError(nil, 400, "bad request: username is missing")
	}
	if user.Password == "" {
		return nil, tools.NewError(nil, 400, "bad request: password is missing")
	}
	if len(user.Password) < handle.C.PasswordMinLength {
		return nil, tools.NewError(nil, 400, "bad request: password is too short")
	}
	if user.Domains == nil || len(user.Domains) == 0 {
		return nil, tools.NewError(nil, 400, "bad request: domains is missing")
	}
	if user.Variables == nil {
		user.Variables = make(map[string]interface{})
	}

	uid, err := user.Create()
	return CreateResponse{
		Status: "ok",
		UserID: uid.Hex(),
	}, err
}
