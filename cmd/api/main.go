package main

import (
	"importerapi/config"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error load .env file")
	}
}

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	// Initialize routes
	bootstrapRoutes(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6969"
	}

	log.Fatal(app.Listen(":" + port))
}
