
package main

import (
  "log"
  "net/http"

  "github.com/julienschmidt/httprouter"

  c "github.com/konek/auth-server/controllers"
  "github.com/konek/auth-server/models"
  "github.com/konek/auth-server/config"
)

func main() {
  router := httprouter.New()
  conf := config.Config()
  root := conf.Root

  router.POST(root + "/user", handler(conf, c.CreateUser))
  router.GET(root + "/user/:uid", handler(conf, c.InfosUser))
  router.PUT(root + "/user/:uid", handler(conf, c.EditUser))
  router.DELETE(root + "/user/:uid", handler(conf, c.DeleteUser))

  router.GET(root + "/session/:token", handler(conf, c.InfosSession))
  router.POST(root + "/session", handler(conf, c.Login))
  router.PUT(root + "/session", handler(conf, c.Check))
  router.DELETE(root + "/session/:token", handler(conf, c.Logoff))

  router.POST(root + "/auth", handler(conf, c.Auth))
  router.GET(root + "/clean/:age", handler(conf, c.Clean))

  router.GET(root + "/list/users", handler(conf, c.ListUsers))
  router.GET(root + "/list/sessions", handler(conf, c.ListSessions))

  err := models.Init(conf)
  if err != nil {
    log.Fatal(err)
  }
  log.Fatal(http.ListenAndServe(conf.Listen, router))
}

