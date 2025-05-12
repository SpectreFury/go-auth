package main

import (
	"log"

	"github.com/SpectreFury/go-auth/internal/env"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error getting environment variables")
		return
	}

	config := config{addr: env.GetString("PORT", "localhost:3000")}

	app := &application{
		config: config,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
