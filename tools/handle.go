package tools

import (
	"github.com/julienschmidt/httprouter"
	"net/http"

	"bitbucket.org/konek/auth-server/config"
)

// Handle is used to reduce the controllers function signature.
type Handle struct {
	R *http.Request
	P httprouter.Params
	C config.Conf
}
