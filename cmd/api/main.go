package main

import (
	"log"

	"github.com/likhon22/social/internal/env"
	"github.com/likhon22/social/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":4000"),
	}
	app := &application{
		config: cfg,
	}
	store := store.NewStorage(nil)
	mux := app.mount()
	log.Fatal(app.serve(mux))

}
