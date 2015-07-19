package controllers

import (
	"bitbucket.org/konek/mgo"
	"gopkg.in/konek/auth-server.v1/models"
	"gopkg.in/konek/auth-server.v1/tools"
)

// ListUsersResponse ...
type ListUsersResponse struct {
	Status string        `json:"status"`
	Users  []models.User `json:"users"`
}

// ListUsers returns a list of all users
func ListUsers(handle tools.Handle, db *mgo.DbQueue) (interface{}, error) {
	list, err := models.ListUsers(db)
	if err != nil {
		return nil, err
	}

	for i := range list {
		list[i].Password = ""
		list[i].Salt = ""
	}

	return ListUsersResponse{
		Status: "ok",
		Users:  list,
	}, nil
}
