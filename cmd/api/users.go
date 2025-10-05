package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"github.com/likhon22/social/internal/store"
)

type contextKey string

const userIDKey contextKey = "userID"

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
func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.StatusInternalServerError(w, r, err)
		return
	}
}

type FollowUser struct {
	UserId int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	follower := getUserFromContext(r)

	var payload FollowUser

	if err := readJSON(w, r, &payload); err != nil {
		app.StatusInternalServerError(w, r, err)
		return

	}
	log.Println(follower.ID, follower.ID, payload.UserId)
	err := app.store.Followers.Follow(r.Context(), payload.UserId, follower.ID)
	if err != nil {
		// Check if it's a Postgres unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				writeJSONError(w, http.StatusConflict, "you already followed")
				return
			}
		}

		app.StatusInternalServerError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, "you followed successfully")
}

func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollower := getUserFromContext(r)
	var payload FollowUser

	if err := readJSON(w, r, &payload); err != nil {
		app.StatusInternalServerError(w, r, err)
		return

	}
	log.Println(unFollower.ID, unFollower.ID, payload.UserId)
	err := app.store.Followers.UnFOllow(r.Context(), payload.UserId, unFollower.ID)
	if err != nil {
		app.StatusInternalServerError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, "you unfollowed successfully")
}

// user middleware
func (app *application) UserIdContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userIdParam := chi.URLParam(r, "userId")
		if userIdParam == "" {
			http.Error(w, "user ID missing", http.StatusBadRequest)
			return
		}
		userId, err := strconv.ParseInt(userIdParam, 10, 64)
		if err != nil {
			app.BadRequestError(w, r, err)
			return
		}
		if userId < 1 {
			app.BadRequestError(w, r, errors.New("invalid ID"))
			return
		}
		user, err := app.store.Users.GetUserById(ctx, userId)
		if err != nil {
			app.StatusInternalServerError(w, r, err)
			return

		}
		ctx = context.WithValue(r.Context(), userIDKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userIDKey).(*store.User)
	return user

}
