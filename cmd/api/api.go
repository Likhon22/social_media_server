package main

import (
	"log"
	"net/http"
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
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Get("/", app.getPostsHandler)

			r.Route("/{postId}", func(r chi.Router) {
				r.Get("/", app.getPostByIDHandler)
				r.Delete("/", app.deletePostByIDHandler)
				r.Patch("/", app.updatePostHandler)
			})

		})

		//users

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)
			r.Get("/", app.getUserHandler)

			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.UserIdContextMiddleware)
				r.Get("/", app.getUserByIdHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unFollowUserHandler)
			})

			r.Get("/email", app.getUserByEmailHandler)

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})

		})
		//comment
		r.Route("/comments", func(r chi.Router) {
			r.Post("/", app.CreateCommentHandler)
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
