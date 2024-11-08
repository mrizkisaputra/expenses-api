package repository

import (
	"context"
	. "github.com/mrizkisaputra/expenses-api/internal/user"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// userPostgresRepository is data/repository implementation
// of service layer UserPostgresRepository
type userPostgresRepository struct {
	db *gorm.DB
}

// NewUserPostgresRepository is a factory for initializing User Repositories
func NewUserPostgresRepository(db *gorm.DB) UserPostgresRepository {
	return &userPostgresRepository{db: db}
}

func (u *userPostgresRepository) Create(ctx context.Context, entity *model.User) (*model.User, error) {
	db := u.db.WithContext(ctx)
	/**
	INSERT INTO "users" ("email","password","first_name","last_name","created_at","updated_at")
	VALUES (?,?,?,?,?,?) RETURNING "id"
	*/
	if err := db.Omit("avatar", "city", "phone_number").Create(entity).Error; err != nil {
		return nil, errors.Wrap(err, "UserPostgresRepository.Register.Create")
	}
	return entity, nil
}

func (u *userPostgresRepository) Update(ctx context.Context, entity *model.User) (*model.User, error) {
	DB := u.db.WithContext(ctx)

	/**
	UPDATE "users" SET "email"=?,"password"=?,"first_name"=?,"last_name"=?,"avatar"=?,"city"=?,"phone_number"=?,"updated_at"=? WHERE id=?
	*/
	if err := DB.Model(model.User{}).Where("id = ?", entity.Id).Updates(entity).Error; err != nil {
		return nil, errors.Wrap(err, "UserPostgresRepository.Update.Updates")
	}
	return entity, nil
}

func (u *userPostgresRepository) FindByEmail(ctx context.Context, entity *model.User) (*model.User, error) {
	DB := u.db.WithContext(ctx)

	// SELECT * FROM "users" WHERE "users"."email" = ? LIMIT 1
	if err := DB.Where(model.User{Email: entity.Email}).Take(entity).Error; err != nil {
		return nil, errors.Wrap(err, "UserPostgresRepository.FindByEmail.Take")
	}
	return entity, nil
}

func (u *userPostgresRepository) FindById(ctx context.Context, entity *model.User) (*model.User, error) {
	DB := u.db.WithContext(ctx)

	// SELECT * FROM "users" WHERE "users"."id" = ? LIMIT 1
	if err := DB.Take(entity).Error; err != nil {
		return nil, errors.Wrap(err, "UserPostgresRepository.FindById.Take")
	}
	return entity, nil
}

func (u *userPostgresRepository) FindAlreadyExistByEmail(ctx context.Context, entity *model.User) (int64, error) {
	var total int64
	DB := u.db.WithContext(ctx)

	// SELECT count(*) FROM "users" WHERE "users"."email" = ?
	if err := DB.Model(model.User{}).Where(model.User{Email: entity.Email}).Count(&total).Error; err != nil {
		return 0, errors.Wrap(err, "UserPostgresRepository.FindAlreadyExistByEmail.Count")
	}
	return total, nil
}
