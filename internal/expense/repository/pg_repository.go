package repository

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/internal/expense"
	"github.com/mrizkisaputra/expenses-api/internal/expense/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewExpensePgRepository(db *gorm.DB) expense.PostgresRepository {
	return &postgresRepository{
		db: db,
	}
}

func (p *postgresRepository) Create(ctx context.Context, expense *model.Expense) error {
	/**
	SQL: INSERT INTO "expenses" ("id_user","description","amount","category","created_at","updated_at","deleted_at")
	VALUES (?,?,?,?,?,?,NULL) RETURNING "id"
	*/
	DB := p.db.WithContext(ctx)
	if err := DB.Create(expense).Error; err != nil {
		return errors.Wrap(err, "postgresRepository.Create")
	}
	return nil
}

func (p *postgresRepository) FindByIdAndUserId(ctx context.Context, expense *model.Expense, id, userId string) error {
	/**
	SQL: SELECT * FROM "expenses" WHERE (id = '?' AND id_user = '?') AND "expenses"."deleted_at" IS NULL LIMIT 1
	*/
	DB := p.db.WithContext(ctx)
	if err := DB.Where("id = ? AND id_user = ?", id, userId).Take(expense).Error; err != nil {
		return errors.Wrap(err, "postgresRepository.FindById")
	}
	return nil
}

func (p *postgresRepository) Remove(ctx context.Context, expense *model.Expense) error {
	/**
	SQL: UPDATE "expenses" SET "deleted_at"='2024-11-21 21:25:19.08' WHERE "expenses"."id" = '?' AND "expenses"."deleted_at" IS NULL
	*/
	DB := p.db.WithContext(ctx)
	if err := DB.Delete(expense).Error; err != nil {
		return errors.Wrap(err, "postgresRepository.Remove")
	}
	return nil
}

func (p *postgresRepository) Update(ctx context.Context, expense *model.Expense) error {
	/**
	SQL: UPDATE "expenses"
	SET "id_user"='?',"description"='?',"amount"=?,"category"='?',"updated_at"=? WHERE (id = '?' AND id_user = '?') AND "expenses"."deleted_at" IS NULL AND "id" = '?'
	*/
	DB := p.db.WithContext(ctx)
	if err := DB.Where("id = ? AND id_user = ?", expense.Id, expense.UserId).Updates(expense).Error; err != nil {
		return errors.Wrap(err, "postgresRepository.Update")
	}
	return nil
}
