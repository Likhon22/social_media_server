package main

import (
	"log"

	"github.com/likhon22/social/internal/config"
	"github.com/likhon22/social/internal/db"
	"github.com/likhon22/social/internal/store"

	"github.com/likhon22/social/internal/env"
)

func main() {
	dbConfig := &config.DbConfig{
		Addr:         env.GetString("DB_ADDR", "postgres://username:password@localhost/social?sslmode=disable"),
		MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
		MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}
	conn, err := db.NewDB(*dbConfig)
	store := store.NewStorage(conn)
	if err != nil {
		log.Fatal(err)
	}
	db.Seed(*store)

}
