package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/likhon22/social/internal/config"
)

var dbInstance *sql.DB

func NewDB(cfg config.DbConfig) (*sql.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}
	db, err := sql.Open("postgres", cfg.Addr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	dbInstance = db
	return dbInstance, nil

}
