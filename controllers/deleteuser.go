package controllers

import (
	"bitbucket.org/konek/mgo"
	"gopkg.in/konek/auth-server.v1/models"
	"gopkg.in/konek/auth-server.v1/tools"
)

// DeleteResponse ...
type DeleteResponse struct {
	Status string `json:"status"`
}

// DeleteUser the user, this does not clean its sessions, yet (TODO).
func DeleteUser(handle tools.Handle, db *mgo.DbQueue) (interface{}, error) {
	var user models.User

	uid := handle.P.ByName("uid")
	if uid == "" {
		return nil, tools.NewError(nil, 400, "bad request: missing userID")
	}
	if len(uid) != 24 {
		return nil, tools.NewError(nil, 400, "bad request: invalid userID")
	}

	user.IDFromHex(uid)
	err := user.Delete(db)
	if err != nil {
		return nil, err
	}
	return DeleteResponse{"ok"}, nil
}
