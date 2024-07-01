// /repository/expense_repository.go

package repository

import (
	"database/sql"
	"errors"
	"time"

	"livecode-catatan-keuangan/models"
)

type ExpenseRepository interface {
	CreateExpense(expense models.Expense) (models.Expense, error)
	GetLastBalance(userID string) (float64, error)
	ListExpenses(userID string, page, size int) ([]models.Expense, int, error)
	GetExpenseByID(userID, id string) (models.Expense, error)
	GetExpensesByType(userID, transactionType string) ([]models.Expense, error)
}

type expenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &expenseRepository{db}
}

func (r *expenseRepository) CreateExpense(expense models.Expense) (models.Expense, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Expense{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var balance float64
	query := `SELECT balance FROM expenses WHERE user_id = $1 ORDER BY date DESC LIMIT 1`
	err = tx.QueryRow(query, expense.UserID).Scan(&balance)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return models.Expense{}, err
	}

	if expense.TransactionType == "CREDIT" {
		balance += expense.Amount
	} else if expense.TransactionType == "DEBIT" {
		if balance < expense.Amount {
			tx.Rollback()
			return models.Expense{}, errors.New("insufficient balance")
		}
		balance -= expense.Amount
	}
	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()
	expense.Balance = balance

	insertQuery := `INSERT INTO expenses ( date, amount, transaction_type, description, created_at, updated_at, user_id, balance) 
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	err = tx.QueryRow(insertQuery, expense.Date, expense.Amount, expense.TransactionType, expense.Description, expense.CreatedAt, expense.UpdatedAt, expense.UserID, balance).Scan(
		&expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt, &expense.UserID, &expense.Balance)
	if err != nil {
		tx.Rollback()
		return models.Expense{}, err
	}

	err = tx.Commit()
	if err != nil {
		return models.Expense{}, err
	}

	return expense, nil
}

func (r *expenseRepository) GetLastBalance(userID string) (float64, error) {
	var balance float64
	query := `SELECT balance FROM expenses WHERE user_id = $1 ORDER BY date DESC LIMIT 1`
	err := r.db.QueryRow(query, userID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return balance, err
}

func (r *expenseRepository) ListExpenses(userID string, page, size int) ([]models.Expense, int, error) {
	offset := (page - 1) * size
	query := `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE user_id = $1 LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, userID, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Balance, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, expense)
	}

	countQuery := `SELECT COUNT(*) FROM expenses WHERE user_id = $1`
	var totalRows int
	err = r.db.QueryRow(countQuery, userID).Scan(&totalRows)
	if err != nil {
		return nil, 0, err
	}

	return expenses, totalRows, nil
}

func (r *expenseRepository) GetExpenseByID(userID, id string) (models.Expense, error) {
	var expense models.Expense
	query := `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE user_id = $1 AND id = $2`
	err := r.db.QueryRow(query, userID, id).Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Balance, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt)
	if err == sql.ErrNoRows {
		return expense, errors.New("expense not found")
	}
	return expense, err
}

func (r *expenseRepository) GetExpensesByType(userID, transactionType string) ([]models.Expense, error) {
	query := `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE user_id = $1 AND transaction_type = $2`
	rows, err := r.db.Query(query, userID, transactionType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Balance, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}
