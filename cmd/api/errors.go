package main

import (
	"net/http"
)

func (app *application) StatusInternalServerError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("internal server ", "error", err.Error(), "method", r.Method, "URL", r.URL.Path)
	writeJSONError(w, http.StatusInternalServerError, "something happen on the server")
}
func (app *application) BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	var msg string
	if err != nil {
		msg = err.Error()
	} else {
		msg = "bad request"
	}

	app.logger.Errorw("bad request", "error", err.Error(), "method", r.Method, "URL", r.URL.Path)
	writeJSONError(w, http.StatusBadRequest, msg)
}
func (app *application) NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	var msg string
	if err != nil {
		msg = err.Error()
	} else {
		msg = "not found"
	}

	app.logger.Errorw("not found", "error", msg, "method", r.Method, "URL", r.URL.Path)
	writeJSONError(w, http.StatusNotFound, "not found")
}
