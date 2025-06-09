package main

import (
	"importerapi/config"
	"importerapi/internal/services"
	"importerapi/internal/worker"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("NODE_ENV") != "development" {
		log.Println("Running in production mode, skipping .env loading")
		return
	}
	// Load environment variables from .env file
	// This is only for development purposes
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
		BodyLimit: 20 * 1024 * 1024, // 20 MB
	})

	// Job queue for import jobs
	jobQueue := make(chan worker.ImportJob, 10)

	// setup import service
	importService := services.NewImportService(db)

	// Initialize worker service
	workerService := worker.Worker{
		Service: importService,
	}

	// Start worker pool with 3 workers
	for i := range 3 {
		go workerService.StartImportWorker(jobQueue, i)
	}

	// Initialize routes
	bootstrapRoutes(app, db, jobQueue)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6969"
	}

	log.Fatal(app.Listen(":" + port))
}
