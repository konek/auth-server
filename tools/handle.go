package tools

import (
	"net/http"

	"go.konek.io/auth-server/config"
	"go.konek.io/rest"
)

// Handle is used to reduce the controllers function signature.
type Handle struct {
	R *http.Request
	P rest.Params
	C config.Conf
}
