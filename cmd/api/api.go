package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware

	"github.com/temideewan/go-social/docs" // required to generate the swagger docs
	"github.com/temideewan/go-social/internal/store"
)

type application struct {
	config   config
	store    store.Storage
	errorLog *log.Logger
	infoLog  *log.Logger
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiUrl string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		docsUrl := fmt.Sprintf("%s/v1/swagger/doc.json", app.config.apiUrl)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Get("/", app.getAllPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Put("/", app.updatePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)

			})
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
	})
	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = fmt.Sprintf("%s/", strings.TrimPrefix(app.config.apiUrl, "http://"))
	docs.SwaggerInfo.BasePath = "v1"
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	app.infoLog.Printf("Server has started at %s: ", app.config.addr)
	return srv.ListenAndServe()
}
