package main

import (
	"log"
	"os"

	"github.com/temideewan/go-social/internal/db"
	"github.com/temideewan/go-social/internal/env"
	"github.com/temideewan/go-social/internal/store"

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
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
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

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		errorLog.Panic(err)
	}

	defer db.Close()
	infoLog.Println("Database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config:   cfg,
		store:    store,
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := app.mount()
	errorLog.Fatal(app.run(mux))
}
