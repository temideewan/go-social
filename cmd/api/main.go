package main

import (
	"github.com/temideewan/go-social/internal/db"
	"github.com/temideewan/go-social/internal/env"
	"github.com/temideewan/go-social/internal/store"
	"go.uber.org/zap"

	_ "github.com/temideewan/go-social/docs"
)

const version = "0.0.1"

//	@title			Td Gopher social
//	@description	API for TDGopher a social network I'm building
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	APiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				The api assigns a eky when you sign up. You need to pass it in the "Authorization" header for endpoints that require authentication.

func main() {
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		env:    env.GetString("ENV", "development"),
		apiUrl: env.GetString("EXTERNAL_URL", "http://localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	// database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
