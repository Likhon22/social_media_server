package main

import (
	"log"
	"net/http"
)

func (app *application) StatusInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path %s method: %s", err.Error(), r.Method, r.URL.Path)
	writeJSONError(w, http.StatusInternalServerError, "something happen on the server")
}
func (app *application) BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	var msg string
	if err != nil {
		msg = err.Error()
	} else {
		msg = "bad request"
	}
	log.Printf("bad request error: %s path %s method: %s", msg, r.Method, r.URL.Path)
	writeJSONError(w, http.StatusInternalServerError, msg)
}
