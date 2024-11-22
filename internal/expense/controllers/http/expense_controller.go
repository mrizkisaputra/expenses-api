package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/mrizkisaputra/expenses-api/internal/expense"
	"github.com/mrizkisaputra/expenses-api/internal/expense/model"
	"github.com/mrizkisaputra/expenses-api/internal/middleware"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	"github.com/mrizkisaputra/expenses-api/pkg/utils"
	"github.com/mrizkisaputra/expenses-api/pkg/validator"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ControllerConfig struct {
	ExpenseService expense.ExpenseService
	Logger         *logrus.Logger
}

// expenseController acts as a struct for injecting an implementation of ExpenseController interface
// for use in controller methods
type expenseController struct {
	expenseService expense.ExpenseService
	logger         *logrus.Logger
}

// NewExpenseController is a factory function
// initializing a expenseController with its service layer dependencies
func NewExpenseController(config *ControllerConfig) expense.ExpenseController {
	return &expenseController{
		expenseService: config.ExpenseService,
		logger:         config.Logger,
	}
}

func (ec expenseController) CreateNewExpense() gin.HandlerFunc {
	type CreateExpenseRequest struct {
		Description string   `json:"description" validate:"required,max=200"`
		Amount      *float64 `json:"amount" validate:"required,gt=0"`
		Category    string   `json:"category" validate:"required,max=100"`
	}

	return func(ctx *gin.Context) {
		auth := middleware.GetAuth(ctx)

		request := new(CreateExpenseRequest)
		if err := utils.ReadRequest(ctx, request, binding.JSON); err != nil {
			utils.LogErrorResponse(ctx, ec.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		entity := &model.Expense{
			UserId:      auth.Id,
			Description: request.Description,
			Amount:      request.Amount,
			Category:    request.Category,
		}
		response, err := ec.expenseService.Insert(ctx, entity)
		if err != nil {
			utils.LogErrorResponse(ctx, ec.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		ctx.JSON(http.StatusCreated, &model.ApiResponse{
			Status:  http.StatusCreated,
			Message: "Created",
			Data:    response,
		})
	}
}

func (ec expenseController) GetExpenseById() gin.HandlerFunc {
	type GetExpenseRequest struct {
		Id     string `validate:"required,uuid"`
		UserId string
	}

	return func(ctx *gin.Context) {
		auth := middleware.GetAuth(ctx)

		request := &GetExpenseRequest{
			Id:     ctx.Param("id"),
			UserId: auth.Id.String(),
		}

		if err := validator.ValidateStruct(ctx, request); err != nil {
			utils.LogErrorResponse(ctx, ec.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		response, err := ec.expenseService.GetById(ctx, request.Id, request.UserId)
		if err != nil {
			utils.LogErrorResponse(ctx, ec.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		ctx.JSON(http.StatusOK, &model.ApiResponse{
			Status:  http.StatusOK,
			Message: "OK",
			Data:    response,
		})
	}
}

func (ec expenseController) DeleteExpense() gin.HandlerFunc {
	type DeleteExpenseRequest struct {
		Id string `validate:"required,uuid"`
	}

	return func(ctx *gin.Context) {
		auth := middleware.GetAuth(ctx)
		request := &DeleteExpenseRequest{
			Id: ctx.Param("id"),
		}

		if err := validator.ValidateStruct(ctx, request); err != nil {
			utils.LogErrorResponse(ctx, ec.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		entity := &model.Expense{
			Id:     uuid.MustParse(request.Id),
			UserId: auth.Id,
		}
		if err := ec.expenseService.Delete(ctx, entity); err != nil {
			utils.LogErrorResponse(ctx, ec.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		ctx.JSON(http.StatusNoContent, &model.ApiResponse{
			Status:  http.StatusNoContent,
			Message: "No Content",
		})
	}
}

func (ec expenseController) UpdateExpense() gin.HandlerFunc {
	type UpdateExpenseRequest struct {
		Id          string   `json:"-" validate:"required,uuid"`
		Description string   `json:"description" validate:"omitempty,max=200"`
		Amount      *float64 `json:"amount" validate:"omitempty,gt=0"`
		Category    string   `json:"category" validate:"omitempty,max=100"`
	}

	return func(ctx *gin.Context) {
		auth := middleware.GetAuth(ctx)

		request := new(UpdateExpenseRequest)
		request.Id = ctx.Param("id")
		if err := utils.ReadRequest(ctx, request, binding.JSON); err != nil {
			utils.LogErrorResponse(ctx, ec.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		entity := &model.Expense{
			Id:          uuid.MustParse(request.Id),
			UserId:      auth.Id,
			Description: request.Description,
			Amount:      request.Amount,
			Category:    request.Category,
		}
		response, err := ec.expenseService.Update(ctx, entity)
		if err != nil {
			utils.LogErrorResponse(ctx, ec.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		ctx.JSON(http.StatusAccepted, &model.ApiResponse{
			Status:  http.StatusAccepted,
			Message: "Accepted",
			Data:    response,
		})
	}
}

func (ec expenseController) GetAllExpense() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
