package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SpectreFury/go-auth/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	conn   *pgx.Conn
	ctx    context.Context
}

type config struct {
	addr string
}

func (app *application) mount() *chi.Mux {
	// Connect to the database before creating router
	ctx := context.Background()

	conn, err := db.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	app.conn = conn
	app.ctx = ctx

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", app.healthHandler)
		r.Get("/session", app.sessionHandler)
		r.Post("/signup", app.signUpHandler)
		r.Post("/login", app.loginHandler)
		r.Get("/logout", app.logoutHandler)
	})

	return r
}

func (app *application) run(mux *chi.Mux) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("Server listening on %s", app.config.addr)

	return server.ListenAndServe()
}
