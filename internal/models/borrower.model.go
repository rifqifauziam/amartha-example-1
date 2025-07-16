package models

import (
	"time"

	"gorm.io/gorm"
)

// Borrower represents a borrower entity
type Borrower struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FirstName string         `gorm:"not null" json:"first_name"`
	LastName  string         `gorm:"not null" json:"last_name"`
	Email     string         `gorm:"not null;unique" json:"email"`
	Phone     string         `gorm:"not null" json:"phone"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Loans     []Loan         `gorm:"foreignKey:BorrowerID" json:"loans,omitempty"`
}
