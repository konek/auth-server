package controllers

import (
	"go.konek.io/auth-server/models"
	"go.konek.io/auth-server/tools"
)

// ListSessionsResponse ...
type ListSessionsResponse struct {
	Status   string           `json:"status"`
	Sessions []models.Session `json:"sessions"`
}

// ListSessions returns a list of all the sessions (expired or not)
func ListSessions(handle tools.Handle) (interface{}, error) {
	list, err := models.ListSessions()
	if err != nil {
		return nil, err
	}
	return ListSessionsResponse{
		Status:   "ok",
		Sessions: list,
	}, nil
}
