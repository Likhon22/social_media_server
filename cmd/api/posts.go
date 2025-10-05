package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/likhon22/social/internal/store"
)

type CreatePostPayload struct {
	Content string   `json:"content" db:"content" validator:"required,max=1000"`
	Title   string   `json:"title" db:"title" validator:"required,max=100"`
	Tags    []string `json:"tags" db:"tags" validate:"required"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.BadRequestError(w, r, err)
	}
	if err := Validate.Struct(payload); err != nil {
		app.BadRequestError(w, r, err)
	}
	posts := &store.Post{}

	if err := readJSON(w, r, posts); err != nil {
		app.BadRequestError(w, r, err)
		return
	}

	err := app.store.Posts.Create(r.Context(), posts)
	if err != nil {
		log.Println(err)
		app.StatusInternalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusCreated, posts); err != nil {
		app.StatusInternalServerError(w, r, err)
		return
	}

}

func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {

	posts, err := app.store.Posts.GetAll(r.Context())
	if err != nil {
		log.Println(err)
		app.BadRequestError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, posts); err != nil {
		app.StatusInternalServerError(w, r, err)
		return
	}

}

func (app *application) getPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	postIDParam := chi.URLParam(r, "postId")
	if postIDParam == "" {
		app.BadRequestError(w, r, errors.New("postID is needed"))
		return
	}
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		app.BadRequestError(w, r, err)
		return
	}
	if postID < 1 {
		app.BadRequestError(w, r, errors.New("invalid ID"))
		return
	}

	post, err := app.store.Posts.GetByID(r.Context(), postID)
	if err != nil {
		log.Println(err)
		app.StatusInternalServerError(w, r, err)
		return
	}
	if post == nil {
		app.NotFoundError(w, r, err)
		return
	}
	comments, err := app.store.Comments.GetCommentsWithPost(r.Context(), postID)
	post.Comments = *comments
	if err != nil {
		app.StatusInternalServerError(w, r, err)

	}
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.StatusInternalServerError(w, r, err)
	}
}
