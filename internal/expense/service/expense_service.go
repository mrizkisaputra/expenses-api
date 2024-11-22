package service

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/mrizkisaputra/expenses-api/internal/expense"
	"github.com/mrizkisaputra/expenses-api/internal/expense/model"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	"github.com/sirupsen/logrus"
)

type ServiceConfig struct {
	PgRepo expense.PostgresRepository
	Config *config.Config
	Logger *logrus.Logger
}

type expenseService struct {
	pgRepo expense.PostgresRepository
}

func NewExpenseService(config *ServiceConfig) expense.ExpenseService {
	return &expenseService{
		pgRepo: config.PgRepo,
	}
}

func (e *expenseService) Insert(ctx context.Context, request *model.Expense) (*model.Expense, error) {
	request.PrepareCreate()
	if err := e.pgRepo.Create(ctx, request); err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}
	return request, nil
}

func (e *expenseService) Delete(ctx context.Context, request *model.Expense) error {
	expenses := new(model.Expense)
	if err := e.pgRepo.FindByIdAndUserId(ctx, expenses, request.Id.String(), request.UserId.String()); err != nil {
		return httpErrors.NewNotFoundError(err)
	}

	if err := e.pgRepo.Remove(ctx, expenses); err != nil {
		return httpErrors.NewInternalServerError(err)
	}
	return nil
}

func (e *expenseService) Update(ctx context.Context, request *model.Expense) (*model.Expense, error) {
	expenses := new(model.Expense)
	if err := e.pgRepo.FindByIdAndUserId(ctx, expenses, request.Id.String(), request.UserId.String()); err != nil {
		return nil, httpErrors.NewNotFoundError(err)
	}

	request.PrepareUpdate(expenses)

	if err := e.pgRepo.Update(ctx, expenses); err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}
	return expenses, nil
}

func (e *expenseService) GetById(ctx context.Context, id, userId string) (*model.Expense, error) {
	expenses := new(model.Expense)
	if err := e.pgRepo.FindByIdAndUserId(ctx, expenses, id, userId); err != nil {
		return nil, httpErrors.NewNotFoundError(err)
	}
	return expenses, nil
}
