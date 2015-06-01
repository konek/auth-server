package controllers

import (
	"go.konek.io/auth-server/models"
	"go.konek.io/auth-server/tools"
	"go.konek.io/mgo"
	"go.konek.io/rest"

	"github.com/asaskevich/govalidator"
)

// CreateResponse ...
type CreateResponse struct {
	Status string `json:"status"`
	UserID string `json:"uid"`
}

// CreateUser create a new user. Checks for duplicate users and password-length requirement
func CreateUser(handle tools.Handle, db *mgo.DbQueue) (interface{}, error) {
	var user models.User

	user.Enable = true
	user.Domains = nil
	user.Variables = nil
	err := rest.Parse(handle.R, &user)
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
	if govalidator.IsEmail(user.Username) == false {
		return nil, tools.NewError(nil, 400, "bad request: username must be a valid email")
	}

	user.Username, err = govalidator.NormalizeEmail(user.Username)
	if err != nil {
		return nil, tools.NewError(nil, 400, "bad request: username must be a valid email")
	}
	uid, err := user.Create(db)
	return CreateResponse{
		Status: "ok",
		UserID: uid.Hex(),
	}, err
}
