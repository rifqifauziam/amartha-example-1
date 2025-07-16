package services

import (
	"errors"
	"fmt"
	"time"

	"AmarthaExample1/internal/dto"
	"AmarthaExample1/internal/models"
	"AmarthaExample1/internal/repositories"
)

// LoanService handles business logic for loans
type LoanService struct {
	repo *repositories.LoanRepository
}

// NewLoanService creates a new loan service instance
func NewLoanService(repo *repositories.LoanRepository) *LoanService {
	return &LoanService{repo: repo}
}

// CreateLoan creates a new loan with payment schedule
func (s *LoanService) CreateLoan(borrowerID uint, amount float64) (*models.Loan, error) {
	interestRate := 0.10
	totalWeeks := 50

	totalAmount := amount * (1 + interestRate)
	weeklyPayment := totalAmount / float64(totalWeeks)

	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 7*totalWeeks)

	loan := &models.Loan{
		BorrowerID:    borrowerID,
		Amount:        amount,
		InterestRate:  interestRate,
		TotalAmount:   totalAmount,
		WeeklyPayment: weeklyPayment,
		TotalWeeks:    totalWeeks,
		StartDate:     startDate,
		EndDate:       endDate,
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(loan); err != nil {
		return nil, err
	}

	return loan, nil
}

// GetLoanByID retrieves a loan by its ID
func (s *LoanService) GetLoanByID(id uint) (*models.Loan, error) {
	return s.repo.GetByID(id)
}

// GetOutstanding returns the current outstanding amount on a loan
func (s *LoanService) GetOutstanding(loanID uint) (float64, error) {
	return s.repo.GetOutstandingAmount(loanID)
}

// IsDelinquent checks if a borrower is delinquent (missed 2+ consecutive payments)
func (s *LoanService) IsDelinquent(loanID uint) (bool, error) {
	missedCount, err := s.repo.GetMissedPaymentsCount(loanID)
	if err != nil {
		return false, err
	}

	return missedCount >= 2, nil
}

// MakePayment processes a payment for a loan
func (s *LoanService) MakePayment(loanID uint, amount float64) error {
	loan, err := s.repo.GetByID(loanID)
	if err != nil {
		return err
	}

	payments, err := s.repo.GetPaymentsByLoanID(loanID)
	if err != nil {
		return err
	}

	// Check for delinquency
	missedCount, err := s.repo.GetMissedPaymentsCount(loanID)
	if err != nil {
		return err
	}

	// Count pending payments that need to be paid
	var pendingPayments []*models.Payment
	for i := range payments {
		if payments[i].Status == "pending" {
			pendingPayments = append(pendingPayments, &payments[i])
			// Only collect the first payment if not delinquent, or collect all missed payments if delinquent
			if missedCount < 2 && len(pendingPayments) == 1 {
				break
			}
			// If delinquent, collect all missed payments up to the missed count
			if len(pendingPayments) == missedCount {
				break
			}
		}
	}

	if len(pendingPayments) == 0 {
		return errors.New("no pending payments found")
	}

	// Calculate required payment amount based on delinquency
	requiredAmount := loan.WeeklyPayment
	if missedCount >= 2 {
		requiredAmount = loan.WeeklyPayment * float64(len(pendingPayments))
	}

	if amount != requiredAmount {
		if missedCount >= 2 {
			return fmt.Errorf("delinquent loan: payment amount must be %v for %d missed payments", requiredAmount, len(pendingPayments))
		}
		return errors.New("payment amount must match the weekly payment amount")
	}

	// Process payment for the first pending payment or all missed payments if delinquent
	now := time.Now()
	
	if missedCount >= 2 {
		// Pay all missed payments if delinquent
		for _, payment := range pendingPayments {
			payment.Status = "paid"
			payment.PaidDate = &now
			payment.PaymentDate = &now
			payment.UpdatedAt = now

			if err := s.repo.UpdatePayment(payment); err != nil {
				return err
			}
		}
	} else {
		// Pay just the first pending payment if not delinquent
		pendingPayments[0].Status = "paid"
		pendingPayments[0].PaidDate = &now
		pendingPayments[0].PaymentDate = &now
		pendingPayments[0].UpdatedAt = now

		if err := s.repo.UpdatePayment(pendingPayments[0]); err != nil {
			return err
		}
	}

	// Check if all payments are now paid
	allPaid := true
	for _, p := range payments {
		if p.Status != "paid" {
			allPaid = false
			break
		}
	}

	if allPaid {
		loan.Status = "completed"
		loan.UpdatedAt = now
		return s.repo.UpdateLoan(loan)
	}

	return nil
}

// GetLoanSchedule returns the payment schedule for a loan
func (s *LoanService) GetLoanSchedule(loanID uint) ([]dto.ScheduleItemDTO, error) {
	return s.repo.GetLoanSchedule(loanID)
}

// GetPaymentsByLoanID returns all payments for a loan
func (s *LoanService) GetPaymentsByLoanID(loanID uint) ([]models.Payment, error) {
	return s.repo.GetPaymentsByLoanID(loanID)
}
