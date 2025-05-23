# Go Auth Learning Project

This project was created as a personal learning exercise to understand how authentication works in Go. It demonstrates a simple authentication system using Go's `net/http` package, PostgreSQL for data storage, and minimal dependencies.

## Features

- User signup with hashed password storage
- User login with session management (using cookies)
- Protected routes that require authentication
- Logout functionality

## Tech Stack

- **Backend:** Go, net/http, [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt), PostgreSQL
- **Frontend:** React + TypeScript + Vite (see [`web/`](web/README.md))

## Why I Made This

I built this project to learn and experiment with authentication flows in Go, including password hashing, session management, and secure cookie handling. The code is intentionally simple and meant for educational purposes only.

## Running the Project

1. **Clone the repository**
2. **Set up your `.env` file** with the required environment variables (see `.env.example` if available)
3. **Start the backend server:**
   ```sh
   go run ./cmd/api
   ```
4. **Start the frontend:**
   ```sh
   cd web
   npm install
   npm run dev
   ```

## Disclaimer

This project is for learning and demonstration purposes only. Do not use it as-is in production.

---

Made with ❤️ by [SpectreFury](https://github.com/SpectreFury)
