package expense

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/internal/expense/model"
)

// PostgresRepository defines methods the service layer expects.
// any repositories it interacts with to implement
type PostgresRepository interface {
	Create(ctx context.Context, expense *model.Expense) error

	FindByIdAndUserId(ctx context.Context, expense *model.Expense, id, userId string) error

	Remove(ctx context.Context, expense *model.Expense) error

	Update(ctx context.Context, expense *model.Expense) error

	FindAll(ctx context.Context, userId string, offset, limit int) ([]model.Expense, int64, error)

	FindAllByDateRange(ctx context.Context, userId string, startDate, endDate int64, offset, limit int) ([]model.Expense, int64, error)
}
