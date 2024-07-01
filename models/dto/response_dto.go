package dto

type Status struct {
	ResponseCode int `json:"responseCode"`
}

type SingleResponse struct {
	Status Status `json:"status"`
	Data   any    `json:"data"`
}

type Paging struct {
	Page      int `json:"page"`
	TotalData int `json:"totalData"`
}

type PagingResponse struct {
	Status Status `json:"status"`
	Data   []any  `json:"data"`
	Paging Paging `json:"paging"`
}

type CreateExpenseDTO struct {
	Amount          float64 `json:"amount" binding:"required"`
	TransactionType string  `json:"transactionType" binding:"required"`
	Description     string  `json:"description" binding:"required"`
}
type ExpenseResponseDTO struct {
	ID              string  `json:"id"`
	Date            string  `json:"date"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transactionType"`
	Balance         float64 `json:"balance"`
	Description     string  `json:"description"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

type PagedExpenseResponseDTO struct {
	ResponseCode string               `json:"responseCode"`
	Data         []ExpenseResponseDTO `json:"data"`
	Paging       PagingResponseDTO    `json:"paging"`
}

type PagingResponseDTO struct {
	Page      int `json:"page"`
	TotalData int `json:"totalData"`
}
