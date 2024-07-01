package models

import (
	"time"
)

type Expense struct {
	ID              string    `json:"id"`
	Date            time.Time `json:"date"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	Balance         float64   `json:"balance"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ExpenseDTO struct {
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	Description     string  `json:"description"`
}
