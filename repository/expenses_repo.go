package repository

import (
	"database/sql"
	"livecode-catatan-keuangan/models"
)

type ExpenseRepository interface {
	Create(expense *models.Expense) error
	FindAll(page, size int, startDate, endDate string) ([]models.Expense, int, error)
	FindByID(id string) (*models.Expense, error)
	FindByType(transactionType string, page, size int) ([]models.Expense, int, error)
	GetLatestExpense() (*models.Expense, error)
}

type expenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &expenseRepository{db}
}

func (r *expenseRepository) Create(expense *models.Expense) error {
	query := `INSERT INTO expenses (date, amount, transaction_type, balance, description, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, expense.Date, expense.Amount, expense.TransactionType, expense.Balance, expense.Description, expense.CreatedAt, expense.UpdatedAt)
	return err
}

func (r *expenseRepository) FindAll(page, size int, startDate, endDate string) ([]models.Expense, int, error) {
	var rows *sql.Rows
	var err error
	if startDate != "" && endDate != "" {
		query := `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at
		          FROM expenses
		          WHERE date BETWEEN $1 AND $2
		          ORDER BY date DESC
		          LIMIT $3 OFFSET $4`
		rows, err = r.db.Query(query, startDate, endDate, size, (page-1)*size)
	} else {
		query := `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at
		          FROM expenses
		          ORDER BY date DESC
		          LIMIT $1 OFFSET $2`
		rows, err = r.db.Query(query, size, (page-1)*size)
	}

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		expense := models.Expense{}
		err := rows.Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Balance, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, expense)
	}

	countQuery := `SELECT COUNT(*) FROM expenses`
	var totalData int
	err = r.db.QueryRow(countQuery).Scan(&totalData)
	if err != nil {
		return nil, 0, err
	}

	return expenses, totalData, nil
}

func (r *expenseRepository) FindByID(id string) (*models.Expense, error) {
	expense := new(models.Expense)
	query := `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at
	          FROM expenses
	          WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Balance, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (r *expenseRepository) FindByType(transactionType string, page, size int) ([]models.Expense, int, error) {
	query := `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at
	          FROM expenses
	          WHERE transaction_type = $1
	          ORDER BY date DESC
	          LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, transactionType, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		expense := models.Expense{}
		err := rows.Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Balance, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, expense)
	}

	countQuery := `SELECT COUNT(*) FROM expenses WHERE transaction_type = $1`
	var totalData int
	err = r.db.QueryRow(countQuery, transactionType).Scan(&totalData)
	if err != nil {
		return nil, 0, err
	}

	return expenses, totalData, nil
}

func (r *expenseRepository) GetLatestExpense() (*models.Expense, error) {
	query := `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at
	          FROM expenses
	          ORDER BY date DESC, created_at DESC
	          LIMIT 1`
	row := r.db.QueryRow(query)
	expense := new(models.Expense)
	err := row.Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Balance, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return expense, nil
}
