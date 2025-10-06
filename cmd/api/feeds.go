package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	//pagination,filters
	feed, err := app.store.Posts.GetUserFeed(r.Context(), int64(42))
	if err != nil {
		app.StatusInternalServerError(w, r, err)
	}
	if err := writeJSON(w, http.StatusOK, feed); err != nil {
		app.StatusInternalServerError(w, r, err)
	}
}
