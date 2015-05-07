package controllers

import (
	"go.konek.io/auth-server/models"
	"go.konek.io/auth-server/tools"
	"go.konek.io/mgo"
)

// InfosUserResponse ...
type InfosUserResponse struct {
	Status string      `json:"status"`
	Infos  models.User `json:"infos"`
}

// InfosUser returns informations about the user. It does not lists its sessions, yet (TODO)
func InfosUser(handle tools.Handle, db *mgo.DbQueue) (interface{}, error) {
	var user models.User

	uid := handle.P.ByName("uid")
	if handle.C.Public == true {
		ret, err := CheckSession(CheckRequest{
			Domain: "/io/konek/app/user",
			Token:  handle.Sid,
		}, db)
		if err != nil {
			return nil, err
		}
		sess := ret.(CheckResponse)
		if uid == "" {
			uid = sess.Session.UserID
		}
		if sess.Session.UserID != uid {
			return nil, tools.NewError(nil, 403, "forbiden: this is not your account")
		}
	}
	if tools.CheckID(uid) == false {
		return nil, tools.NewError(nil, 400, "bad request: invalid userID")
	}

	user.IDFromHex(uid)
	err := user.Get(db)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return InfosUserResponse{
		Status: "ok",
		Infos:  user,
	}, nil
}
