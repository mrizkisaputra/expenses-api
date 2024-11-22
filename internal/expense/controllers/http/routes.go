package http

import (
	"github.com/gin-gonic/gin"
	. "github.com/mrizkisaputra/expenses-api/internal/expense"
	"github.com/mrizkisaputra/expenses-api/internal/middleware"
)

func MapExpenseRoutes(authGroup *gin.RouterGroup, controller ExpenseController, mw *middleware.MiddlewareManager) {
	authGroup.Use(mw.AuthJwtMiddleware())
	authGroup.POST("/create", controller.CreateNewExpense())
	authGroup.GET("/:id", controller.GetExpenseById())
	authGroup.DELETE("/:id", controller.DeleteExpense())
	authGroup.PATCH("/:id", controller.UpdateExpense())
}
