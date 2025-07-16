package main

import (
	"log"
	"os"

	"AmarthaExample1/internal/config"
	"AmarthaExample1/internal/handlers"
	"AmarthaExample1/internal/models"
	"AmarthaExample1/internal/repositories"
	"AmarthaExample1/internal/routes"
	"AmarthaExample1/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Database configuration
	dbConfig := config.DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "billing_engine"),
	}
	db := config.GetDBInstance(dbConfig)
	defer db.Close()
	db.Conn.AutoMigrate(&models.Loan{}, &models.Payment{}, &models.Borrower{})

	// Initialize repositories
	loanRepo := repositories.NewLoanRepository(db.Conn)

	// Initialize services
	loanService := services.NewLoanService(loanRepo)

	// Initialize handlers
	loanHandler := handlers.NewLoanHandler(loanService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(recover.New())

	routes.SetupLoanRoutes(app, loanHandler)

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
