
package controllers

import (
  "github.com/konek/auth-server/tools"
  "github.com/konek/auth-server/models"
)

// LoginRequest ...
type LoginRequest struct {
  Domain    string `json:"domain"`
  Username  string `json:"username"`
  Password  string `json:"password"`
}

// LoginResponse ...
type LoginResponse struct {
  Status    string `json:"status"`
  Session   SessionResponse `json:"session"`
}

// Login a user, creating a new session.
func Login(handle tools.Handle) (interface{}, error) {
  var q LoginRequest
  var user models.User
  var session models.Session
  var resp LoginResponse

  err := tools.ParseBody(handle.R.Body, &q)
  if err != nil {
    return nil, tools.NewError(err, 400, "bad request: couldn't parse body")
  }

  if q.Domain == "" {
    return nil, tools.NewError(nil, 400, "bad request: domain is missing")
  }
  if q.Domain == "/" {
    return nil, tools.NewError(nil, 400, "bad request: illegal domain")
  }

  if q.Username == "" {
    return nil, tools.NewError(nil, 400, "bad request: username is missing")
  }
  if q.Password == "" {
    return nil, tools.NewError(nil, 400, "bad request: password is missing")
  }

  user.Username = q.Username
  user.Password = q.Password

  ok, err := user.Check()
  if err != nil {
    return nil, err
  }
  if ok == false {
    return nil, tools.NewError(nil, 403, "forbidden: invalid user or password")
  }
  if user.Enable == false {
    return nil, tools.NewError(nil, 403, "forbidden: user is diabled")
  }
  ok = user.CheckDomain(q.Domain)
  if ok == false {
    return nil, tools.NewError(nil, 403, "forbidden: restricted domain")
  }

  session.UserID = user.ID
  session.Domain = q.Domain
  remaining, err := session.Create(handle.C.SessionLifespan)
  if err != nil {
    return nil, err
  }
  resp.Status = "ok"
  resp.Session.Token = session.ID.Hex()
  resp.Session.UserID = session.UserID.Hex()
  resp.Session.Expire = session.Expire
  resp.Session.Remaining = remaining
  return resp, nil
}
