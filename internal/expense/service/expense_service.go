package service

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/mrizkisaputra/expenses-api/internal/expense"
	"github.com/mrizkisaputra/expenses-api/internal/expense/model"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
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

func (e *expenseService) GetAll(ctx context.Context, userId string, request *model.SearchExpenseRequestQueryParam) ([]model.Expense, int64, error) {
	var start, end time.Time
	var now = time.Now()

	switch request.Filter {
	case "last_week":
		{
			start = now.AddDate(0, 0, -7)
			end = now
			break
		}
	case "last_month":
		{
			start = now.AddDate(0, -1, 0)
			end = now
			break
		}
	case "last_3_month":
		{
			start = now.AddDate(0, -3, 0)
			end = now
			break
		}
	case "custom":
		{
			startDate, err := time.Parse("2006-01-02", request.StartDate)
			if err != nil {
				return nil, 0, httpErrors.NewInternalServerError(errors.Wrap(err, "expenseService.GetAll.time.Parse"))
			}

			endDate, err := time.Parse("2006-01-02", request.EndDate)
			if err != nil {
				return nil, 0, httpErrors.NewInternalServerError(errors.Wrap(err, "expenseService.GetAll.time.Parse"))
			}

			start = startDate
			end = endDate
			break
		}
	default:
		expenses, total, err := e.pgRepo.FindAll(ctx, userId, request.Page, request.Limit)
		if err != nil {
			return nil, 0, httpErrors.NewInternalServerError(err)
		}
		return expenses, total, nil
	}

	expenses, total, err := e.pgRepo.FindAllByDateRange(
		ctx,
		userId,
		start.UnixMilli(),
		end.UnixMilli(),
		request.Page,
		request.Limit,
	)
	if err != nil {
		return nil, 0, httpErrors.NewInternalServerError(err)
	}

	return expenses, total, nil
}
