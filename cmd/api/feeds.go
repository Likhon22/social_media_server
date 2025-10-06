package main

import (
	"net/http"

	"github.com/likhon22/social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	//pagination,filters
	fq := store.PaginatedFeedQuery{
		Page:  1,
		Limit: 3,
		Sort:  "desc",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		app.BadRequestError(w, r, err)

	}
	if err := Validate.Struct(fq); err != nil {
		app.BadRequestError(w, r, err)
	}
	feed, err := app.store.Posts.GetUserFeed(r.Context(), int64(42), fq)
	if err != nil {
		app.StatusInternalServerError(w, r, err)
	}
	if err := writeJSON(w, http.StatusOK, feed); err != nil {
		app.StatusInternalServerError(w, r, err)
	}
}
