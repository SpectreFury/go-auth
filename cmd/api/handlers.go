package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SpectreFury/go-auth/internal/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Healthy"}`))
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
	userId, hashedPassword, err := db.GetUser(app.conn, app.ctx, `SELECT id, hashed_password FROM users WHERE email = $1`, user.Email)
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

	// Create a cookie for the user and save it in a sessions table

	sessionId := uuid.New()

	// Save it in database
	err = db.InsertSession(app.conn, app.ctx, `INSERT INTO sessions (session_id, user_id) VALUES ($1, $2)`, sessionId.String(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the cookie to the frontend
	cookie := &http.Cookie{
		Name:     "sessionId",
		Value:    sessionId.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	// User is successfully logged in
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"message":"Login was successful"}`))

}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the cookie from the request
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Delete the session from the database
	err = db.DeleteSession(app.conn, app.ctx, `DELETE FROM sessions WHERE session_id = $1`, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Delete the cookie from the browser
	c := &http.Cookie{
		Name:     "sessionId",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, c)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Logout was successful"}`))
}

func (app *application) protectedHandler(w http.ResponseWriter, r *http.Request) {
	// Get the cookie from the request
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the session exists in the database
	session, err := db.SessionExists(app.conn, app.ctx, `SELECT session_id FROM sessions WHERE session_id = $1`, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !session {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// If the session exists, return a protected resource

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Protected route"}`))
}
