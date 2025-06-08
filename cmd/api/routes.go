package main

import (
	"importerapi/internal/handlers"
	"importerapi/internal/middleware"
	"importerapi/internal/worker"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func bootstrapRoutes(app *fiber.App, db *gorm.DB, jobQueue chan worker.ImportJob) {
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("FRONTEND_URL"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "API is running",
		})
	})

	importHandler := handlers.NewImportHandler(jobQueue, db)
	customerHandler := handlers.NewCustomerHandler(db)
	invoiceHandler := handlers.NewInvoiceHandler(db)
	invoiceItemHandler := handlers.NewInvoiceItemHandler(db)
	productHandler := handlers.NewProductHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(db)

	/* Auth Routes */
	api.Post("/register", authHandler.RegisterHandler)
	api.Post("/login", authHandler.LoginHandler)

	/* Middleware for authentication */
	app.Use(middleware.AuthMiddleware)

	/* Import Routes */
	api.Post("/import/xml", importHandler.ImportXMLDataHandler)
	api.Get("/import/:id/status", importHandler.GetImportStatusHandler)

	/* Customer Routes */
	api.Get("/customers", customerHandler.GetCustomersHandler)

	/* Invoice Routes */
	api.Get("/invoices", invoiceHandler.GetInvoicesHandler)

	/* Invoice Item Routes */
	api.Get("/invoice-items", invoiceItemHandler.GetInvoiceItemsHandler)
	api.Get("/invoice-items/:id", invoiceItemHandler.GetInvoiceItemByID)

	/* Product Routes */
	api.Get("/products", productHandler.GetProducts)
	api.Get("/products/:id", productHandler.GetProductByID)

	/* Summary Route */
	api.Get("/summary", invoiceItemHandler.GetInvoiceItemSummaryHandler)

	/* User Routes */
	api.Get("/user", userHandler.GetProfileHanlder)
}
