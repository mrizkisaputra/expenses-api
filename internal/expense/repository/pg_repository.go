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

func (p *postgresRepository) FindAll(ctx context.Context, userId string, offset, limit int) ([]model.Expense, int64, error) {
	DB := p.db.WithContext(ctx)

	/**
	SQL: SELECT * FROM "expenses" WHERE id_user = ? AND "expenses"."deleted_at" IS NULL LIMIT 10
	*/
	var expenses []model.Expense
	if err := DB.Offset((offset-1)*limit).
		Limit(limit).
		Where("id_user = ?", userId).
		Find(&expenses).Error; err != nil {
		return nil, 0, errors.Wrap(err, "postgresRepository.FindAll.Find")
	}

	/**
	SQL: SELECT count(*) FROM "expenses" WHERE "expenses"."deleted_at" IS NULL
	*/
	var count int64
	if err := DB.Model(&model.Expense{}).Count(&count).Error; err != nil {
		return nil, 0, errors.Wrap(err, "postgresRepository.FindAll.Count")
	}

	return expenses, count, nil
}

func (p *postgresRepository) FindAllByDateRange(
	ctx context.Context,
	userId string,
	startDate, endDate int64,
	offset, limit int,
) ([]model.Expense, int64, error) {
	DB := p.db.WithContext(ctx)

	/**
	SQL:
	*/
	var expenses []model.Expense
	if err := DB.Where("id_user = ? AND created_at BETWEEN ? AND ?", userId, startDate, endDate).
		Order("created_at DESC").
		Offset((offset - 1) * limit).
		Limit(limit).
		Find(&expenses).Error; err != nil {

		return nil, 0, errors.Wrap(err, "postgresRepository.FindAllByDateRange")
	}

	/**
	SQL:
	*/
	var total int64
	if err := DB.Model(&model.Expense{}).
		Where("id_user = ? AND created_at BETWEEN ? AND ?", userId, startDate, endDate).
		Count(&total).Error; err != nil {

		return nil, 0, errors.Wrap(err, "postgresRepository.FindAllByDateRange.Count")
	}

	return expenses, total, nil
}
