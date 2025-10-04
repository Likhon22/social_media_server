package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/likhon22/social/internal/config"
	"github.com/likhon22/social/internal/store"
)

type application struct {
	Config *config.AppConfig
	store  *store.Storage
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.HandleFunc("GET /health", app.healthCheckHandler)
		r.Route("/post", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Get("/", app.getPostsHandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Get("/", app.getPostByIDHandler)
			})

		})

		//users

		r.Route("/user", func(r chi.Router) {
			r.Post("/", app.createUserHandler)

		})
	})
	return r
}
func (app *application) serve(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server on %s", app.Config.Addr)
	return srv.ListenAndServe()
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

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
		app.BadRequestError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.StatusInternalServerError(w, r, err)
	}
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

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
