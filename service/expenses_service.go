package service

import (
	"errors"
	"livecode-catatan-keuangan/models"
	"livecode-catatan-keuangan/models/dto"
	"livecode-catatan-keuangan/repository"
	"time"
)

type ExpenseService interface {
	CreateExpense(createExpenseDTO *dto.CreateExpenseDTO) (*models.Expense, error)
	ListExpenses(page, size int, startDate, endDate string) ([]models.Expense, int, error)
	GetExpenseByID(id string) (*models.Expense, error)
	GetExpensesByType(transactionType string, page, size int) ([]models.Expense, int, error)
}

type expenseService struct {
	expenseRepo repository.ExpenseRepository
}

func NewExpenseService(expenseRepo repository.ExpenseRepository) ExpenseService {
	return &expenseService{expenseRepo}
}

func (s *expenseService) CreateExpense(createExpenseDTO *dto.CreateExpenseDTO) (*models.Expense, error) {
	// Get the latest balance
	latestExpense, err := s.expenseRepo.GetLatestExpense()
	if err != nil {
		return nil, err
	}

	// Calculate new balance
	var balance float64
	if latestExpense != nil {
		balance = latestExpense.Balance
	} else {
		balance = 0
	}

	switch createExpenseDTO.TransactionType {
	case "CREDIT":
		balance += createExpenseDTO.Amount
	case "DEBIT":
		balance -= createExpenseDTO.Amount
	default:
		return nil, errors.New("invalid transaction type")
	}

	expense := &models.Expense{
		Date:            time.Now(),
		Amount:          createExpenseDTO.Amount,
		TransactionType: createExpenseDTO.TransactionType,
		Balance:         balance,
		Description:     createExpenseDTO.Description,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.expenseRepo.Create(expense); err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *expenseService) ListExpenses(page, size int, startDate, endDate string) ([]models.Expense, int, error) {
	return s.expenseRepo.FindAll(page, size, startDate, endDate)
}

func (s *expenseService) GetExpenseByID(id string) (*models.Expense, error) {
	return s.expenseRepo.FindByID(id)
}

func (s *expenseService) GetExpensesByType(transactionType string, page, size int) ([]models.Expense, int, error) {
	return s.expenseRepo.FindByType(transactionType, page, size)
}
