package controllers

import (
	"strconv"

	"go.konek.io/auth-server/models"
	"go.konek.io/auth-server/tools"
	"go.konek.io/mgo"
)

// CleanResponse ...
type CleanResponse struct {
	Status  string `json:"status"`
	Deleted int    `json:"deleted"`
}

// Clean every expired sessions older than age
func Clean(handle tools.Handle, db *mgo.DbQueue) (interface{}, error) {
	sAge := handle.P.ByName("age")

	if sAge == "" {
		return nil, tools.NewError(nil, 400, "bad request: age is missing")
	}
	age, err := strconv.ParseInt(sAge, 10, 64)
	if err != nil {
		return nil, tools.NewError(err, 400, "bad request: invalid age")
	}
	n, err := models.CleanSessions(db, age)
	if err != nil {
		return nil, err
	}
	return CleanResponse{
		Status:  "ok",
		Deleted: n,
	}, nil
}
