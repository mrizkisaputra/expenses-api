package mock

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/internal/expense/model"
	"github.com/stretchr/testify/mock"
)

type MockPostgresRepository struct {
	mock.Mock
}

func (m MockPostgresRepository) Create(ctx context.Context, expense *model.Expense) error {
	//TODO implement me
	panic("implement me")
}

func (m MockPostgresRepository) FindByIdAndUserId(ctx context.Context, expense *model.Expense, id, userId string) error {
	//TODO implement me
	panic("implement me")
}

func (m MockPostgresRepository) Remove(ctx context.Context, expense *model.Expense) error {
	//TODO implement me
	panic("implement me")
}

func (m MockPostgresRepository) Update(ctx context.Context, expense *model.Expense) error {
	//TODO implement me
	panic("implement me")
}

func (m MockPostgresRepository) FindAll(ctx context.Context, userId string, offset, limit int) ([]model.Expense, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockPostgresRepository) FindAllByDateRange(ctx context.Context, userId string, startDate, endDate int64, offset, limit int) ([]model.Expense, int64, error) {
	//TODO implement me
	panic("implement me")
}
