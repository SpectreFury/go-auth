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

func SessionExists(conn *pgx.Conn, ctx context.Context, query string, sessionId string) (bool, error) {
	var scannedSesssionId string

	err := conn.QueryRow(ctx, query, sessionId).Scan(&scannedSesssionId)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Println("RESULT: No session found")
			return false, nil
		}

		fmt.Println("ERROR during query: ", err)
		return false, err
	}

	fmt.Println("Session found: ", scannedSesssionId)
	return true, nil
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

func GetUser(conn *pgx.Conn, ctx context.Context, query string, email string) (userId string, hashedPassword string, returnError error) {
	var scannedId string
	var scannedPassword string

	err := conn.QueryRow(ctx, query, email).Scan(&scannedId, &scannedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Println("RESULT: No user found")
			return "", "", nil
		}

		fmt.Println("ERROR during query: ", err)
		return "", "", err
	}

	fmt.Println("User found: ", scannedId)
	return scannedId, scannedPassword, nil
}

func InsertUser(conn *pgx.Conn, ctx context.Context, query string, email string, hashedPassword string) error {
	_, err := conn.Exec(ctx, query, email, hashedPassword)

	return err
}

func InsertSession(conn *pgx.Conn, ctx context.Context, query string, sessionId string, userId string) error {
	_, err := conn.Exec(ctx, query, sessionId, userId)
	if err != nil {
		fmt.Println("Error in inserting session")
		fmt.Println(err)
		return err
	}

	fmt.Println("Successfully inserted the session")
	return nil
}

func DeleteSession(conn *pgx.Conn, ctx context.Context, query string, sessionId string) error {
	_, err := conn.Exec(ctx, query, sessionId)

	if err != nil {
		fmt.Println("Error in deleting session")
		return err
	}

	fmt.Println("Successfully deleted the session")
	return nil
}
