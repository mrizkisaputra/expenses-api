package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

// Expense mapping table 'expenses'
type Expense struct {
	Id          uuid.UUID      `gorm:"column:id;primary_key;default:uuid_generate_v4();<-:create"` // allow read and create
	UserId      uuid.UUID      `gorm:"column:id_user"`
	Description string         `gorm:"column:description"`
	Amount      *float64       `gorm:"column:amount"`
	Category    string         `gorm:"column:category"`
	CreatedAt   int64          `gorm:"column:created_at;autoCreateTime;<-:create"` // allow read and create
	UpdatedAt   int64          `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
	//User        model.User     `gorm:"foreignKey:user_id;references:id"`
}

func (expense *Expense) TableName() string {
	return "expenses"
}

func (expense *Expense) PrepareCreate() {
	expense.Description = strings.ToLower(strings.TrimSpace(expense.Description))
	expense.Category = strings.ToLower(strings.TrimSpace(expense.Category))
}

func (expense *Expense) PrepareUpdate(oldExpense *Expense) {
	if expense.Category != "" {
		oldExpense.Category = expense.Category
	}

	if expense.Description != "" {
		oldExpense.Description = expense.Description
	}

	if expense.Amount != nil {
		oldExpense.Amount = expense.Amount
	}
}
