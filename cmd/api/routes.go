package main

import (
	"importerapi/internal/handlers"
	"importerapi/internal/worker"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func bootstrapRoutes(app *fiber.App, db *gorm.DB, jobQueue chan worker.ImportJob) {
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "API is running",
		})
	})

	importHandler := handlers.NewImportHandler(jobQueue)
	customerHandler := handlers.NewCustomerHandler(db)
	invoiceHandler := handlers.NewInvoiceHandler(db)
	productHandler := handlers.NewProductHandler(db)

	/* Import Routes */
	api.Post("/import/xml", importHandler.ImportXMLDataHandler)

	/* Customer Routes */
	api.Get("/customers", customerHandler.GetCustomersHandler)

	/* Invoice Routes */
	api.Get("/invoices", invoiceHandler.GetInvoicesHandler)

	/* Product Routes */
	api.Get("/products", productHandler.GetProducts)
}
