package models

import "time"

type Expense struct {
	ID              string    `json:"id"`
	Date            time.Time `json:"date"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transactionType"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	UserID          string    `json:"userId"`
	Balance         float64   `json:"balance"`
}
