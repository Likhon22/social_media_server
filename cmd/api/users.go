package main

import (
	"net/http"

	"github.com/likhon22/social/internal/store"
)

type CreateUserPayload struct {
	Username string `json:"username" db:"username" validator:"required,max=100"`
	Password string `json:"password" db:"password" validator:"required,max=22"`
	Email    string `json:"email" db:"email" validator:"required,max=30"`
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
	}
	if err := Validate.Struct(payload); err != nil {
		app.BadRequestError(w, r, err)
	}
	user := &store.User{}
	err := readJSON(w, r, user)
	if err != nil {
		app.BadRequestError(w, r, err)
	}
	err = app.store.Users.Create(r.Context(), user)
	if err != nil {
		app.StatusInternalServerError(w, r, err)
	}
	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.StatusInternalServerError(w, r, err)
	}

}
