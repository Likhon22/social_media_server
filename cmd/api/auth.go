package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
	"github.com/likhon22/social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" db:"username" validator:"required,max=100"`
	Password string `json:"password" db:"password" validator:"required,min=5,max=200"`
	Email    string `json:"email" db:"email" validator:"required,email,max=50"`
}

// @Summary		Register a new user
// @Description	Registers a new user in the system
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			user	body		RegisterUserPayload	true	"User information"
// @Success		201		{object}	store.User			"User registered"
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Router			/users [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
	}
	if err := Validate.Struct(payload); err != nil {
		app.BadRequestError(w, r, err)
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	// hash password
	if err := user.Password.Set(payload.Password); err != nil {
		app.StatusInternalServerError(w, r, err)
		return
	}
	ctx := r.Context()
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	if err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.Config.Mail.Exp); err != nil {
		app.StatusInternalServerError(w, r, err)
	}

	// mail

	if err := writeJSON(w, http.StatusCreated, user); err != nil {
		app.StatusInternalServerError(w, r, err)
	}

}
