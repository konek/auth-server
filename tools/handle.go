package tools

import (
	"github.com/julienschmidt/httprouter"
	"net/http"

	"go.konek.io/auth-server/config"
)

// Handle is used to reduce the controllers function signature.
type Handle struct {
	R *http.Request
	P httprouter.Params
	C config.Conf
}
