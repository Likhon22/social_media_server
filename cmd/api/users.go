package main

import (
	"errors"
	"log"
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
	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	err := app.store.Users.Create(r.Context(), user)
	if err != nil {
		app.StatusInternalServerError(w, r, err)
	}
	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.StatusInternalServerError(w, r, err)
	}

}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.store.Users.GetUsers(r.Context())
	if err != nil {
		log.Println(err)
		app.BadRequestError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, users); err != nil {
		app.StatusInternalServerError(w, r, err)
		return
	}
}

func (app *application) getUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	userEmailParam := r.URL.Query().Get("email")

	if userEmailParam == "" {
		app.BadRequestError(w, r, errors.New("email is needed"))

	}
	user, err := app.store.Users.GetUserByEmail(r.Context(), userEmailParam)
	if err != nil {
		app.StatusInternalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.StatusInternalServerError(w, r, err)
		return
	}
}
