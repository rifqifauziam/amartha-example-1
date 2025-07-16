package models

import (
	"time"

	"gorm.io/gorm"
)

// Loan represents a loan entity
type Loan struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	BorrowerID    uint           `gorm:"not null" json:"borrower_id"`
	Amount        float64        `gorm:"not null" json:"amount"`
	InterestRate  float64        `gorm:"not null" json:"interest_rate"`
	TotalAmount   float64        `gorm:"not null" json:"total_amount"`
	WeeklyPayment float64        `gorm:"not null" json:"weekly_payment"`
	TotalWeeks    int            `gorm:"not null" json:"total_weeks"`
	StartDate     time.Time      `gorm:"not null" json:"start_date"`
	EndDate       time.Time      `gorm:"not null" json:"end_date"`
	Status        string         `gorm:"not null;default:'active'" json:"status"` // active, completed, defaulted
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Payments      []Payment      `gorm:"foreignKey:LoanID" json:"payments,omitempty"`
}

// Payment represents a payment made for a loan
type Payment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	LoanID      uint           `gorm:"not null" json:"loan_id"`
	Amount      float64        `gorm:"not null" json:"amount"`
	WeekNum     int            `gorm:"not null" json:"week_num"`
	DueDate     time.Time      `gorm:"not null" json:"due_date"`
	PaidDate    *time.Time     `json:"paid_date"`
	PaymentDate *time.Time     `json:"payment_date"`
	Status      string         `gorm:"not null;default:'pending'" json:"status"` // pending, paid, missed
	CreatedAt   time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

