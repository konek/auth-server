package tools

import (
	"net/http"

	"gopkg.in/konek/auth-server.v1/config"
	"go.konek.io/rest"
)

// Handle is used to reduce the controllers function signature.
type Handle struct {
	R *http.Request
	P rest.Params
	C config.Conf
	Sid string
}
