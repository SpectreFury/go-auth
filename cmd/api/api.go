package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SpectreFury/go-auth/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	conn   *pgx.Conn
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

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", app.healthCheckHandler)

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
