package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("POSTGRES_URI"))
	if err != nil {
		return conn, err
	}

	return conn, err
}

func UserExists(conn *pgx.Conn, ctx context.Context, query string, email string) (bool, error) {
	var scannedEmail string

	err := conn.QueryRow(ctx, query, email).Scan(&scannedEmail)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Println("RESULT: No user found")
			return false, nil
		}

		fmt.Println("ERROR during query: ", err)
		return false, err
	}

	fmt.Println("User found: ", scannedEmail)
	return true, nil
}

func GetUser(conn *pgx.Conn, ctx context.Context, query string, email string) (userEmail string, hashedPassword string, returnError error) {
	var scannedEmail string
	var scannedPassword string

	err := conn.QueryRow(ctx, query, email).Scan(&scannedEmail, &scannedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Println("RESULT: No user found")
			return "", "", nil
		}

		fmt.Println("ERROR during query: ", err)
		return "", "", err
	}

	fmt.Println("User found: ", scannedEmail)
	return scannedEmail, scannedPassword, nil
}

func InsertUser(conn *pgx.Conn, ctx context.Context, query string, email string, hashedPassword string) error {
	_, err := conn.Exec(ctx, query, email, hashedPassword)

	return err
}
