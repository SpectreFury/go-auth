package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SpectreFury/go-auth/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) signUpHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the body from the request
	var user UserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("Email: ", user.Email)
	fmt.Println("Hashed Password: ", string(hashedPassword))

	// Check if the user already exists
	userExists, err := db.UserExists(app.conn, app.ctx, `SELECT email FROM users where email = $1`, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userExists {
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Save the user into the database
	err = db.InsertUser(app.conn, app.ctx, `INSERT INTO users (email, hashed_password) VALUES ($1, $2)`, user.Email, string(hashedPassword))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"message":"Sign up was successful"}`))
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserRequest

	// Get the request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Find the user with user's email
	_, hashedPassword, err := db.GetUser(app.conn, app.ctx, `SELECT email, hashed_password FROM users WHERE email = $1`, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		fmt.Println("Incorrect passwords")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// User is successfully logged in

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"message":"Login was successful"}`))

}
