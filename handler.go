package main

import (
	"net/http"

	"gopkg.in/konek/mgo.v1"
	"gopkg.in/konek/auth-server.v1/config"
	"gopkg.in/konek/auth-server.v1/controllers"
	"gopkg.in/konek/auth-server.v1/tools"
	"go.konek.io/rest"
)

func handler(conf config.Conf, db *mgo.DbQueue, fn controllers.ControllerFunc) rest.Controller {
	return func(r *http.Request, p rest.Params) (interface{}, error) {
		sid := ""
		for _, c := range r.Cookies() {
			if c.Name == "auth" {
				sid = c.Value
			}
		}
		resp, err := fn(tools.Handle{
			R:   r,
			P:   p,
			C:   conf,
			Sid: sid,
		}, db)
		return resp, err
	}
}
