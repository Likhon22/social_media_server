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

//	@Summary		Create a new user
//	@Description	Registers a new user in the system
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		CreateUserPayload	true	"User information"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users [post]

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

// @Summary		Get all users
// @Description	Retrieves a list of all registered users
// @Tags			Users
// @Produce		json
// @Success		200	{array}		store.User
// @Failure		400	{object}	error
// @Failure		500	{object}	error
// @Router			/users [get]
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

// @Summary		Get a user by email
// @Description	Retrieves a single user by their email
// @Tags			Users
// @Produce		json
// @Param			userEmail	path		int	true	"User Email"
// @Success		200			{object}	store.User
// @Failure		400			{object}	error
// @Failure		404			{object}	error
// @Failure		500			{object}	error
// @Router			/users/{userEmail} [get]
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

// @Summary		Get a user by ID
// @Description	Retrieves a single user by their ID
// @Tags			Users
// @Produce		json
// @Param			userId	path		int	true	"User ID"
// @Success		200		{object}	store.User
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Failure		500		{object}	error
// @Router			/users/{userId} [get]
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

// @Summary		Follow a user
// @Description	Allows the authenticated user to follow another user
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			userId	path		int			true	"User ID to follow"
// @Param			payload	body		FollowUser	true	"Authenticated user info (if needed)"
// @Success		200		{string}	string		"you followed successfully"
// @Failure		400		{object}	error
// @Failure		409		{object}	error	"you already followed"
// @Failure		500		{object}	error
// @Router			/users/{userId}/follow [put]
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

// @Summary		Unfollow a user
// @Description	Allows the authenticated user to unfollow another user
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			userId	path		int			true	"User ID to unfollow"
// @Param			payload	body		FollowUser	false	"Optional payload"
// @Success		200		{string}	string		"you unfollowed successfully"
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Router			/users/{userId}/unfollow [put]
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
