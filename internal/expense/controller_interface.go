package expense

import "github.com/gin-gonic/gin"

type ExpenseController interface {
	CreateNewExpense() gin.HandlerFunc

	GetExpenseById() gin.HandlerFunc

	DeleteExpense() gin.HandlerFunc

	UpdateExpense() gin.HandlerFunc

	GetAllExpense() gin.HandlerFunc
}
