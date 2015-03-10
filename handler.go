package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"bitbucket.org/konek/auth-server/config"
	"bitbucket.org/konek/auth-server/controllers"
	"bitbucket.org/konek/auth-server/tools"
)

func handler(conf config.Conf, fn controllers.ControllerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		resp, err := fn(tools.Handle{
			R: r,
			P: p,
			C: conf,
		})
		if err == nil {
			err = writeJSON(w, resp)
		}
		if err != nil {
			if errDetails, ok := err.(tools.APIError); ok == true {
				writeError(w, errDetails)
			} else {
				log.Printf("Error: %s", err)
				writeError(w, tools.NewError(
					err,
					500,
					"An unexpected error occured, please contact an administrator"))
			}
		}
	}
}

func writeJSON(w http.ResponseWriter, data interface{}) error {
	chunk, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(200)

	_, err = w.Write(chunk)
	return err
}

func writeError(w http.ResponseWriter, e tools.APIError) error {
	chunk, err := json.Marshal(tools.ErrorResponse{
		Status:  "error",
		Message: e.Error(),
	})
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(e.Code)

	_, err = w.Write(chunk)
	return err
}
