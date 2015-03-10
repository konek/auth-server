package controllers

import (
	"github.com/konek/auth-server/models"
	"github.com/konek/auth-server/tools"
)

// DeleteResponse ...
type DeleteResponse struct {
	Status string `json:"status"`
}

// DeleteUser the user, this does not clean its sessions, yet (TODO).
func DeleteUser(handle tools.Handle) (interface{}, error) {
	var user models.User

	uid := handle.P.ByName("uid")
	if uid == "" {
		return nil, tools.NewError(nil, 400, "bad request: missing userID")
	}
	if len(uid) != 24 {
		return nil, tools.NewError(nil, 400, "bad request: invalid userID")
	}

	user.IDFromHex(uid)
	err := user.Delete()
	if err != nil {
		return nil, err
	}
	return DeleteResponse{"ok"}, nil
}
