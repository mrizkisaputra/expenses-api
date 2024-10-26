package model

import "github.com/google/uuid"

// Mapping tabel 'users'
type User struct {
	Id          uuid.UUID `gorm:"column:id;primary_key;default:uuid_generate_v4();<-:create"` // allow read and create
	FirstName   string    `gorm:"column:first_name"`
	LastName    string    `gorm:"column:last_name"`
	Email       string    `gorm:"column:email"`
	Password    string    `gorm:"column:password"`
	Avatar      *string   `gorm:"column:avatar"`
	City        *string   `gorm:"column:city"`
	PhoneNumber *string   `gorm:"column:phone_number"`
	CreatedAt   int64     `gorm:"column:created_at;autoCreateTime:milli;<-:create"` // allow read and create
	UpdatedAt   int64     `gorm:"column:updated_at;autoCreateTime:milli;autoCreateTime:milli"`
}

func (u *User) TableName() string {
	return "users"
}
