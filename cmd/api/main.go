package main

import (
	"log"

	"github.com/likhon22/social/internal/config"
	"github.com/likhon22/social/internal/db"
	"github.com/likhon22/social/internal/env"
	"github.com/likhon22/social/internal/store"
)

func main() {
	cfg := &config.AppConfig{
		Addr: env.GetString("ADDR", ":4000"),
		DB: &config.DbConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://username:password@localhost/social?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	db, err := db.NewDB(*cfg.DB)

	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("Connected to database successfully")
	store := store.NewStorage(db)
	app := &application{
		Config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.serve(mux))

}
