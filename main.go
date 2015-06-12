package main

import (
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/konek/mgo"
	"go.konek.io/auth-server/config"
	c "go.konek.io/auth-server/controllers"
	"go.konek.io/rest"
)

func main() {
	router := rest.New()
	conf := config.Config()
	root := conf.Root

	db := mgo.NewDbQueue(100)
	// Adding 3 connections
	for i := 0; i < 3; i++ {
		db.AddConnection(conf.DbURL, conf.DbName)
	}

	if conf.Public == false {
		fmt.Println("Public = false")
	} else {
		fmt.Println("Public = true")
	}
	router.GET(root+"/user/:uid", handler(conf, db, c.InfosUser))
	router.GET(root+"/user/", handler(conf, db, c.InfosUser)) // UID from session
	if conf.Public == false {
		router.POST(root+"/user", handler(conf, db, c.CreateUser))
		router.PUT(root+"/user/:uid", handler(conf, db, c.EditUser))
		router.DELETE(root+"/user/:uid", handler(conf, db, c.DeleteUser))
	}

	if conf.Public == false {
		router.GET(root+"/session/:token", handler(conf, db, c.InfosSession))
	}
	router.POST(root+"/session", handler(conf, db, c.Login))
	router.PUT(root+"/session", handler(conf, db, c.Check))
	router.DELETE(root+"/session/:token", handler(conf, db, c.Logoff))

	router.POST(root+"/auth", handler(conf, db, c.Auth))
	if conf.Public == false {
		router.GET(root+"/clean/:age", handler(conf, db, c.Clean))
	}

	if conf.Public == false {
		router.GET(root+"/list/users", handler(conf, db, c.ListUsers))
		router.GET(root+"/list/sessions", handler(conf, db, c.ListSessions))
	}

	db.Run()
	fmt.Println("listening on ", conf.Listen)
	log.Fatal(http.ListenAndServe(conf.Listen, router))
}
