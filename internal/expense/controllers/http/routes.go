package http

import (
	"github.com/gin-gonic/gin"
	. "github.com/mrizkisaputra/expenses-api/internal/expense"
	"github.com/mrizkisaputra/expenses-api/internal/middleware"
)

func MapExpenseRoutes(expenseGroup *gin.RouterGroup, controller ExpenseController, mw *middleware.MiddlewareManager) {
	expenseGroup.Use(mw.AuthJwtMiddleware())
	expenseGroup.POST("/create", controller.CreateNewExpense())
	expenseGroup.GET("/:id", controller.GetExpenseById())
	expenseGroup.DELETE("/:id", controller.DeleteExpense())
	expenseGroup.PATCH("/:id", controller.UpdateExpense())
	expenseGroup.GET("/", controller.GetAllExpense())
}
