# Billing Engine

A RESTful API for a loan billing system that provides loan schedules, tracks outstanding amounts, and monitors delinquency status.

## Features

- Loan schedule generation for 50-week loans
- Outstanding amount tracking
- Delinquency status monitoring (borrower is delinquent after 2 consecutive missed payments)
- Payment processing

## Technical Stack

- Go (Golang)
- Fiber v2 (Web Framework)
- GORM (ORM)
- MySQL (Database)
- Docker & Docker Compose

## Project Structure

The project follows a clean architecture pattern with dependency injection:

```
/cmd
  /server
    main.go            # Application entry point
/internal
  /config
    database.go        # Database configuration
  /dto
    loan.dto.go        # Data Transfer Objects
  /handlers
    loan.handler.go    # HTTP Request Handlers
  /models
    loan.model.go      # Database Models
  /repositories
    loan.repository.go # Data Access Layer
  /routes
    loan.route.go      # API Routes
  /services
    loan.service.go    # Business Logic
```

## API Endpoints

- `POST /api/loans` - Create a new loan
- `GET /api/loans/:id` - Get loan details
- `GET /api/loans/:id/outstanding` - Get outstanding amount
- `GET /api/loans/:id/delinquent` - Check if loan is delinquent
- `GET /api/loans/:id/schedule` - Get loan payment schedule
- `POST /api/loans/:id/payment` - Make a payment

## Running the Application

### Using Docker Compose

```bash
docker compose build
docker compose up -d
```

The API will be available at http://localhost:8080

### Running Locally

1. Set up MySQL database
2. Configure environment variables
3. Run the application:

```bash
go run cmd/server/main.go
```

## Loan Terms

- 50-week loan for Rp 5,000,000/-
- Flat interest rate of 10% per annum
- Weekly repayment of Rp 110,000 (total repayment: Rp 5,500,000)
- Borrowers can only pay the exact weekly amount or not pay at all
