package main

import (
	"AmarthaExample1/internal/config"
	"AmarthaExample1/internal/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

func main() {
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "3306")
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "user")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "password")
	}
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "billing_engine")
	}

	dbConfig := config.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	// Retry connection up to 5 times
	var db *config.Database
	for i := 0; i < 5; i++ {
		fmt.Printf("Attempting database connection (attempt %d/5) to %s:%s...\n", i+1, dbConfig.Host, dbConfig.Port)
		db = config.GetDBInstance(dbConfig)
		if db != nil {
			break
		}
		fmt.Println("Connection failed, retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	if db == nil {
		log.Fatal("Failed to get database instance after multiple attempts")
	}

	fmt.Println("Successfully connected to database")

	err := db.Conn.AutoMigrate(&models.Borrower{}, &models.Loan{}, &models.Payment{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	createDummyData(db.Conn)

	fmt.Println("Database setup completed successfully!")
}

func createDummyData(db *gorm.DB) {
	borrower1 := models.Borrower{
		ID:        1,
		FirstName: "Test1",
		LastName:  "Test1",
		Email:     "test1@example.com",
		Phone:     "1234567890",
	}

	if err := db.Create(&borrower1).Error; err != nil {
		log.Fatalf("Failed to create borrower: %v", err)
	}

	borrower2 := models.Borrower{
		ID:        2,
		FirstName: "Test2",
		LastName:  "Test2",
		Email:     "test2@example.com",
		Phone:     "0987654321",
	}

	if err := db.Create(&borrower2).Error; err != nil {
		log.Fatalf("Failed to create borrower: %v", err)
	}

	borrower3 := models.Borrower{
		ID:        3,
		FirstName: "Test3",
		LastName:  "Test3",
		Email:     "test3@example.com",
		Phone:     "5556667777",
	}

	if err := db.Create(&borrower3).Error; err != nil {
		log.Fatalf("Failed to create borrower: %v", err)
	}

	// Loan 1 - normal loan with first 3 weeks paid
	loan1 := models.Loan{
		BorrowerID:    1,
		Amount:        5000000,
		InterestRate:  10,
		TotalAmount:   5500000,
		WeeklyPayment: 110000,
		TotalWeeks:    50,
		StartDate:     time.Now().AddDate(0, 0, -21),
		EndDate:       time.Now().AddDate(0, 0, -21).AddDate(0, 0, 50*7),
		Status:        "active",
	}

	if err := db.Create(&loan1).Error; err != nil {
		log.Fatalf("Failed to create loan: %v", err)
	}

	startDate := loan1.StartDate
	for i := 1; i <= loan1.TotalWeeks; i++ {
		dueDate := startDate.AddDate(0, 0, i*7)
		payment := models.Payment{
			LoanID:  loan1.ID,
			WeekNum: i,
			DueDate: dueDate,
			Amount:  loan1.WeeklyPayment,
			Status:  "pending",
		}

		if i <= 3 {
			paidDate := dueDate
			payment.Status = "paid"
			payment.PaidDate = &paidDate
			payment.PaymentDate = &paidDate
		}

		if err := db.Create(&payment).Error; err != nil {
			log.Fatalf("Failed to create payment: %v", err)
		}
	}

	// Loan 2 - with one missed payment (not delinquent)
	loan2 := models.Loan{
		BorrowerID:    2,
		Amount:        5000000,
		InterestRate:  10,
		TotalAmount:   5500000,
		WeeklyPayment: 110000,
		TotalWeeks:    50,
		StartDate:     time.Now().AddDate(0, 0, -35),
		EndDate:       time.Now().AddDate(0, 0, 315),
		Status:        "active",
	}

	if err := db.Create(&loan2).Error; err != nil {
		log.Fatalf("Failed to create loan: %v", err)
	}

	// Create payment records for all weeks
	for i := 1; i <= loan2.TotalWeeks; i++ {
		dueDate := loan2.StartDate.AddDate(0, 0, i*7)

		// Create payment record (pending by default)
		payment := models.Payment{
			LoanID:  loan2.ID,
			Amount:  loan2.WeeklyPayment,
			WeekNum: i,
			DueDate: dueDate,
			Status:  "pending",
		}

		// For the first 2 weeks, mark as paid
		if i <= 2 {
			paidDate := loan2.StartDate.AddDate(0, 0, i*7)
			payment.Status = "paid"
			payment.PaidDate = &paidDate
			payment.PaymentDate = &paidDate
		}

		if err := db.Create(&payment).Error; err != nil {
			log.Fatalf("Failed to create payment: %v", err)
		}
	}

	// Loan 3 - delinquent loan with multiple consecutive missed payments
	loan3 := models.Loan{
		BorrowerID:    3,
		Amount:        5000000,
		InterestRate:  10,
		TotalAmount:   5500000,
		WeeklyPayment: 110000,
		TotalWeeks:    50,
		StartDate:     time.Now().AddDate(0, 0, -35), // Started 5 weeks ago
		EndDate:       time.Now().AddDate(0, 0, 315),
		Status:        "active",
	}

	if err := db.Create(&loan3).Error; err != nil {
		log.Fatalf("Failed to create loan: %v", err)
	}

	for i := 1; i <= loan3.TotalWeeks; i++ {
		dueDate := loan3.StartDate.AddDate(0, 0, i*7)

		payment := models.Payment{
			LoanID:  loan3.ID,
			WeekNum: i,
			DueDate: dueDate,
			Amount:  loan3.WeeklyPayment,
			Status:  "pending",
		}

		// Mark first 2 payments as paid (weeks 3-5 will be missed)
		if i <= 2 {
			paidDate := dueDate
			payment.Status = "paid"
			payment.PaidDate = &paidDate
			payment.PaymentDate = &paidDate
		}

		if err := db.Create(&payment).Error; err != nil {
			log.Fatalf("Failed to create payment: %v", err)
		}
	}

	fmt.Println("Created dummy loans and payments successfully")
}
