package controllers

import (
	"bitbucket.org/konek/mgo"
	"gopkg.in/konek/auth-server.v1/models"
	"gopkg.in/konek/auth-server.v1/tools"
)

// ListSessionsResponse ...
type ListSessionsResponse struct {
	Status   string           `json:"status"`
	Sessions []models.Session `json:"sessions"`
}

// ListSessions returns a list of all the sessions (expired or not)
func ListSessions(handle tools.Handle, db *mgo.DbQueue) (interface{}, error) {
	list, err := models.ListSessions(db)
	if err != nil {
		return nil, err
	}
	return ListSessionsResponse{
		Status:   "ok",
		Sessions: list,
	}, nil
}
