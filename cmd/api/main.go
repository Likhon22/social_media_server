package main

import (
	"github.com/likhon22/social/internal/config"
	"github.com/likhon22/social/internal/db"
	"github.com/likhon22/social/internal/env"
	"github.com/likhon22/social/internal/store"
	"go.uber.org/zap"
)

//	@title			Swagger GoSocialMedia API
//	@description	API for GoSocialMedia, a social network for people
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	cfg := &config.AppConfig{
		Addr:   env.GetString("ADDR", ":5000"),
		ApiURL: env.GetString("EXTERNAL_URL", "localhost:5000"),
		DB: &config.DbConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://username:password@localhost/social?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Version: env.GetString("VERSION", "0.0.1"),
		Env:     env.GetString("ENV", "development"),
	}
	//logger

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	//database
	db, err := db.NewDB(*cfg.DB)

	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("Connected to database successfully")
	store := store.NewStorage(db)
	app := &application{
		Config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()

	logger.Fatal(app.serve(mux))

}
