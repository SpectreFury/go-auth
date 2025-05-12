package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("POSTGRES_URI"))
	if err != nil {
		return conn, err
	}

	return conn, err
}

func ExecQuery(conn *pgx.Conn, ctx context.Context, query string) error {
	log.Println("Executing query: INSERT")
	_, err := conn.Exec(ctx, query)

	return err
}
