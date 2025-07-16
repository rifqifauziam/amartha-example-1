package handlers

import (
	"AmarthaExample1/internal/dto"
	"AmarthaExample1/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// LoanHandler handles HTTP requests for loans
type LoanHandler struct {
	service *services.LoanService
}

// NewLoanHandler creates a new loan handler instance
func NewLoanHandler(service *services.LoanService) *LoanHandler {
	return &LoanHandler{service: service}
}

// CreateLoan handles the creation of a new loan
func (h *LoanHandler) CreateLoan(c *fiber.Ctx) error {
	var req dto.CreateLoanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	loan, err := h.service.CreateLoan(req.BorrowerID, req.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.LoanResponse{
		ID:            loan.ID,
		BorrowerID:    loan.BorrowerID,
		Amount:        loan.Amount,
		InterestRate:  loan.InterestRate,
		TotalAmount:   loan.TotalAmount,
		WeeklyPayment: loan.WeeklyPayment,
		TotalWeeks:    loan.TotalWeeks,
		StartDate:     loan.StartDate,
		EndDate:       loan.EndDate,
		Status:        loan.Status,
	})
}

// GetLoan handles retrieving a loan by ID
func (h *LoanHandler) GetLoan(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid loan ID",
		})
	}

	loan, err := h.service.GetLoanByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.LoanResponse{
		ID:            loan.ID,
		BorrowerID:    loan.BorrowerID,
		Amount:        loan.Amount,
		InterestRate:  loan.InterestRate,
		TotalAmount:   loan.TotalAmount,
		WeeklyPayment: loan.WeeklyPayment,
		TotalWeeks:    loan.TotalWeeks,
		StartDate:     loan.StartDate,
		EndDate:       loan.EndDate,
		Status:        loan.Status,
	})
}

// GetOutstanding handles retrieving the outstanding amount for a loan
func (h *LoanHandler) GetOutstanding(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid loan ID",
		})
	}

	loan, err := h.service.GetLoanByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	outstanding, err := h.service.GetOutstanding(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.OutstandingResponse{
		LoanID:            loan.ID,
		TotalAmount:       loan.TotalAmount,
		AmountPaid:        loan.TotalAmount - outstanding,
		OutstandingAmount: outstanding,
	})
}

// IsDelinquent handles checking if a loan is delinquent
func (h *LoanHandler) IsDelinquent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid loan ID",
		})
	}

	_, err = h.service.GetLoanByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	isDelinquent, err := h.service.IsDelinquent(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := dto.DelinquencyResponse{
		LoanID:       uint(id),
		IsDelinquent: isDelinquent,
	}

	if isDelinquent {
		response.Reason = "Borrower has missed 2 or more consecutive payments"
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// MakePayment handles processing a payment for a loan
func (h *LoanHandler) MakePayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid loan ID",
		})
	}

	var req dto.PaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Use a different variable to avoid redeclaration
	var loanErr error
	_, loanErr = h.service.GetLoanByID(uint(id))
	if loanErr != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": loanErr.Error(),
		})
	}

	if err := h.service.MakePayment(uint(id), req.Amount); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	outstanding, _ := h.service.GetOutstanding(uint(id))

	return c.Status(fiber.StatusOK).JSON(dto.PaymentResponse{
		Success:      true,
		Message:      "Payment processed successfully",
		RemainingDue: outstanding,
	})
}

// GetLoanSchedule handles retrieving the payment schedule for a loan
func (h *LoanHandler) GetLoanSchedule(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid loan ID",
		})
	}

	_, err = h.service.GetLoanByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	schedule, err := h.service.GetLoanSchedule(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	scheduleItems := make([]dto.ScheduleItemDTO, len(schedule))
	for i, item := range schedule {
		scheduleItems[i] = dto.ScheduleItemDTO{
			WeekNum:     item.WeekNum,
			DueDate:     item.DueDate,
			Amount:      item.Amount,
			Status:      item.Status,
			PaymentDate: item.PaymentDate,
		}
	}

	return c.Status(fiber.StatusOK).JSON(dto.ScheduleResponse{
		LoanID:   uint(id),
		Schedule: scheduleItems,
	})
}
