package controllers

import (
	"go.konek.io/auth-server/models"
	"go.konek.io/auth-server/tools"
)

// InfosUserResponse ...
type InfosUserResponse struct {
	Status string      `json:"status"`
	Infos  models.User `json:"infos"`
}

// InfosUser returns informations about the user. It does not lists its sessions, yet (TODO)
func InfosUser(handle tools.Handle) (interface{}, error) {
	var user models.User

	uid := handle.P.ByName("uid")
	if uid == "" {
		return nil, tools.NewError(nil, 400, "bad request: missing userID")
	}
	if tools.CheckID(uid) == false {
		return nil, tools.NewError(nil, 400, "bad request: invalid userID")
	}

	user.IDFromHex(uid)
	err := user.Get()
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return InfosUserResponse{
		Status: "ok",
		Infos:  user,
	}, nil
}
