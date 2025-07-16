package dto

import "time"

// CreateLoanRequest represents the request to create a new loan
type CreateLoanRequest struct {
	BorrowerID uint    `json:"borrower_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

// LoanResponse represents the loan response
type LoanResponse struct {
	ID            uint      `json:"id"`
	BorrowerID    uint      `json:"borrower_id"`
	Amount        float64   `json:"amount"`
	InterestRate  float64   `json:"interest_rate"`
	TotalAmount   float64   `json:"total_amount"`
	WeeklyPayment float64   `json:"weekly_payment"`
	TotalWeeks    int       `json:"total_weeks"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Status        string    `json:"status"`
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// ScheduleResponse represents the loan schedule response
type ScheduleResponse struct {
	LoanID   uint              `json:"loan_id"`
	Schedule []ScheduleItemDTO `json:"schedule"`
}

// ScheduleItemDTO represents a single item in the loan schedule
type ScheduleItemDTO struct {
	WeekNum     int        `json:"week_num"`
	DueDate     time.Time  `json:"due_date"`
	Amount      float64    `json:"amount"`
	Status      string     `json:"status"`
	PaymentDate *time.Time `json:"payment_date,omitempty"`
}

// OutstandingResponse represents the outstanding amount response
type OutstandingResponse struct {
	LoanID            uint    `json:"loan_id"`
	TotalAmount       float64 `json:"total_amount"`
	AmountPaid        float64 `json:"amount_paid"`
	OutstandingAmount float64 `json:"outstanding_amount"`
}

// DelinquencyResponse represents the delinquency status response
type DelinquencyResponse struct {
	LoanID       uint   `json:"loan_id"`
	IsDelinquent bool   `json:"is_delinquent"`
	Reason       string `json:"reason,omitempty"`
}

// PaymentResponse represents the payment response
type PaymentResponse struct {
	Success      bool    `json:"success"`
	Message      string  `json:"message"`
	PaymentID    uint    `json:"payment_id,omitempty"`
	RemainingDue float64 `json:"remaining_due"`
}
