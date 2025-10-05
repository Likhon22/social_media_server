package main

import (
	"net/http"

	"github.com/likhon22/social/internal/store"
)

type CreateCommentPayload struct {
	PostID  int64  `json:"post_id" db:"post_id" validate:"required"`
	UserID  int64  `json:"user_id" db:"user_id" validate:"required"`
	Content string `json:"content" db:"content" validate:"required,max=1000"`
}

func (app *application) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var commentPayload CreateCommentPayload

	if err := readJSON(w, r, &commentPayload); err != nil {
		app.StatusInternalServerError(w, r, err)
		return
	}
	Validate.Struct(commentPayload)
	comment := &store.Comment{
		PostID:  int(commentPayload.PostID),
		UserID:  int(commentPayload.UserID),
		Content: commentPayload.Content,
	}

	if err := app.store.Comments.CreateComment(r.Context(), comment); err != nil {
		app.StatusInternalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, comment); err != nil {
		app.StatusInternalServerError(w, r, err)
	}

}
