package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	data := map[string]string{"status": "available", "environment": app.Config.Env, "version": app.Config.Version}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		writeJSONError(w, r, http.StatusInternalServerError, err.Error())
	}
}
