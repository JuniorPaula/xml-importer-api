package main

import (
	"importerapi/config"
	"importerapi/internal/worker"
	"log"
	"os"
	"time"

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

	app := fiber.New(fiber.Config{
		BodyLimit:   20 * 1024 * 1024, // 20 MB
		ReadTimeout: 60 * time.Second,
	})

	// Job queue for import jobs
	jobQueue := make(chan worker.ImportJob, 10)
	// Start worker pool with 3 workers
	for i := range 3 {
		go worker.StartImportWorker(jobQueue, i)
	}

	// Initialize routes
	bootstrapRoutes(app, db, jobQueue)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6969"
	}

	log.Fatal(app.Listen(":" + port))
}
