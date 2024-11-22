package expense

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/internal/expense/model"
)

// ExpenseService defines methods the layer controller expects.
// any services it interacts with to implement.
type ExpenseService interface {
	Insert(ctx context.Context, request *model.Expense) (*model.Expense, error)

	Delete(ctx context.Context, request *model.Expense) error

	Update(ctx context.Context, request *model.Expense) (*model.Expense, error)

	GetById(ctx context.Context, id, userId string) (*model.Expense, error)
}
