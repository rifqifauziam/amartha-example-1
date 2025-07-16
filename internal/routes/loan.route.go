package routes

import (
	"AmarthaExample1/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupLoanRoutes sets up all loan related routes
func SetupLoanRoutes(app *fiber.App, handler *handlers.LoanHandler) {
	api := app.Group("/api")
	loans := api.Group("/loans")

	// Loan endpoints
	loans.Post("/", handler.CreateLoan)
	loans.Get("/:id", handler.GetLoan)
	loans.Get("/:id/outstanding", handler.GetOutstanding)
	loans.Get("/:id/delinquent", handler.IsDelinquent)
	loans.Get("/:id/schedule", handler.GetLoanSchedule)
	loans.Post("/:id/payment", handler.MakePayment)
}
