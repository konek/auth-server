package main

import (
	"fmt"
	"net/http"

	"go.konek.io/auth-server/config"
	"go.konek.io/auth-server/controllers"
	"go.konek.io/auth-server/tools"
	"go.konek.io/mgo"
	"go.konek.io/rest"
)

func handler(conf config.Conf, db *mgo.DbQueue, fn controllers.ControllerFunc) rest.Controller {
	return func(r *http.Request, p rest.Params) (interface{}, error) {
		fmt.Println(r)
		resp, err := fn(tools.Handle{
			R: r,
			P: p,
			C: conf,
		}, db)
		return resp, err
	}
}

