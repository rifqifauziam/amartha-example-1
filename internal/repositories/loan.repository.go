package repositories

import (
	"errors"
	"time"

	"AmarthaExample1/internal/dto"
	"AmarthaExample1/internal/models"

	"gorm.io/gorm"
)

// LoanRepository handles database operations for loans
type LoanRepository struct {
	db *gorm.DB
}

// NewLoanRepository creates a new loan repository instance
func NewLoanRepository(db *gorm.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

// Create creates a new loan and its payment schedule
func (r *LoanRepository) Create(loan *models.Loan) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(loan).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create payment schedule
	for i := 0; i < loan.TotalWeeks; i++ {
		dueDate := loan.StartDate.AddDate(0, 0, 7*i)
		payment := models.Payment{
			LoanID:    loan.ID,
			Amount:    loan.WeeklyPayment,
			WeekNum:   i + 1,
			DueDate:   dueDate,
			Status:    "pending",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := tx.Create(&payment).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetByID retrieves a loan by its ID
func (r *LoanRepository) GetByID(id uint) (*models.Loan, error) {
	var loan models.Loan
	if err := r.db.Preload("Payments").First(&loan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("loan not found")
		}
		return nil, err
	}
	return &loan, nil
}

// GetPaymentsByLoanID retrieves all payments for a loan
func (r *LoanRepository) GetPaymentsByLoanID(loanID uint) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Where("loan_id = ?", loanID).Order("week_num").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// UpdatePayment updates a payment record
func (r *LoanRepository) UpdatePayment(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

// GetOutstandingAmount calculates the outstanding amount for a loan
func (r *LoanRepository) GetOutstandingAmount(loanID uint) (float64, error) {
	var totalPaid float64
	if err := r.db.Model(&models.Payment{}).
		Where("loan_id = ? AND status = ?", loanID, "paid").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalPaid).Error; err != nil {
		return 0, err
	}

	var loan models.Loan
	if err := r.db.First(&loan, loanID).Error; err != nil {
		return 0, err
	}

	return loan.TotalAmount - totalPaid, nil
}

// GetMissedPaymentsCount returns the count of consecutive missed payments
func (r *LoanRepository) GetMissedPaymentsCount(loanID uint) (int, error) {
	var payments []models.Payment
	if err := r.db.Where("loan_id = ?", loanID).Order("week_num DESC").Find(&payments).Error; err != nil {
		return 0, err
	}

	consecutiveMissed := 0
	currentTime := time.Now()

	for _, payment := range payments {
		if payment.Status == "pending" && payment.DueDate.Before(currentTime) {
			consecutiveMissed++
		} else if payment.Status == "paid" {
			break
		}
	}

	return consecutiveMissed, nil
}

// GetLoanSchedule returns the complete loan schedule with payment status
func (r *LoanRepository) GetLoanSchedule(loanID uint) ([]dto.ScheduleItemDTO, error) {
	var payments []models.Payment
	if err := r.db.Where("loan_id = ?", loanID).Order("week_num").Find(&payments).Error; err != nil {
		return nil, err
	}

	var schedule []dto.ScheduleItemDTO
	for _, payment := range payments {
		scheduleItem := dto.ScheduleItemDTO{
			WeekNum:   payment.WeekNum,
			DueDate:   payment.DueDate,
			Amount:    payment.Amount,
			Status:    payment.Status,
		}

		if payment.PaidDate != nil {
			scheduleItem.PaymentDate = payment.PaidDate
		}

		schedule = append(schedule, scheduleItem)
	}

	return schedule, nil
}

// UpdateLoan updates a loan record
func (r *LoanRepository) UpdateLoan(loan *models.Loan) error {
	return r.db.Save(loan).Error
}
