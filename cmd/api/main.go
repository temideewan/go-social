package main

import (
	"log"
	"os"

	"github.com/temideewan/go-social/internal/db"
	"github.com/temideewan/go-social/internal/env"
	"github.com/temideewan/go-social/internal/store"
)

const version = "0.0.1"

func main() {
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		env:  env.GetString("ENV", "development"),
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
