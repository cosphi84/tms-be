package main

import (
	"log"
	"tms-be/internal/database"
	"tms-be/internal/offices"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, relying on OS environment variables")
		return
	}
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Perform seeding operations
	offices.Seed(db)

}
