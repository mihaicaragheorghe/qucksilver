package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mihaicaragheorghe/qucksilver/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("**WARNING** Could not load .env file")
	}

	addr := ":" + os.Getenv("PORT")
	server.Start(addr)
}
