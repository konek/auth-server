
package tools

import (
  "net/http"
  "github.com/julienschmidt/httprouter"

  "github.com/konek/auth-server/config"
)

// Handle is used to reduce the controllers function signature.
type Handle struct {
  R *http.Request
  P httprouter.Params
  C config.Conf
}

