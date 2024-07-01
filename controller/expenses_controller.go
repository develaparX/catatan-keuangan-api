package controller

import (
	"livecode-catatan-keuangan/middleware"
	"livecode-catatan-keuangan/models/dto"
	"livecode-catatan-keuangan/service"
	"livecode-catatan-keuangan/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExpenseController struct {
	expenseService service.ExpenseService
	aM             middleware.AuthMiddleware
	rg             *gin.RouterGroup
}

func (c *ExpenseController) CreateExpense(ctx *gin.Context) {
	var createExpenseDTO dto.CreateExpenseDTO
	if err := ctx.ShouldBindJSON(&createExpenseDTO); err != nil {
		utils.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	expense, err := c.expenseService.CreateExpense(&createExpenseDTO)
	if err != nil {
		utils.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSingleResponse(ctx, "Success Create New Expense", expense, http.StatusOK)
}

func (c *ExpenseController) ListExpenses(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	expenses, totalData, err := c.expenseService.ListExpenses(page, size, startDate, endDate)
	if err != nil {
		utils.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	// Type assertion to convert []models.Expense to []any
	var expenseList []any
	for _, expense := range expenses {
		expenseList = append(expenseList, expense)
	}

	utils.SendPagingResponse(ctx, "Success Retrieve Expenses", expenseList, dto.Paging{Page: page, TotalData: totalData}, http.StatusOK)
}

func (c *ExpenseController) GetExpenseByID(ctx *gin.Context) {
	id := ctx.Param("id")
	expense, err := c.expenseService.GetExpenseByID(id)
	if err != nil {
		utils.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSingleResponse(ctx, "Success Retrieve Expense", expense, http.StatusOK)
}

func (c *ExpenseController) GetExpensesByType(ctx *gin.Context) {
	transactionType := ctx.Param("type")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	expenses, totalData, err := c.expenseService.GetExpensesByType(transactionType, page, size)
	if err != nil {
		utils.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	// Type assertion to convert []models.Expense to []any
	var expenseList []any
	for _, expense := range expenses {
		expenseList = append(expenseList, expense)
	}

	utils.SendPagingResponse(ctx, "Success Retrieve Expenses by Type", expenseList, dto.Paging{Page: page, TotalData: totalData}, http.StatusOK)
}

func (e *ExpenseController) Route() {
	router := e.rg.Group("/expenses", e.aM.CheckToken())
	{
		router.POST("/", e.CreateExpense)
		router.GET("/", e.ListExpenses)
		router.GET("/:id", e.GetExpenseByID)
		router.GET("/type/:type", e.GetExpensesByType)
	}
}

func NewExpenseController(eS service.ExpenseService, rg *gin.RouterGroup, am middleware.AuthMiddleware) *ExpenseController {
	return &ExpenseController{expenseService: eS, rg: rg, aM: am}
}
