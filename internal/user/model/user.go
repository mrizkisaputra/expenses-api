package model

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// Mapping tabel 'users'
type User struct {
	Id          uuid.UUID   `gorm:"column:id;primary_key;default:uuid_generate_v4();<-:create"` // allow read and create
	Email       string      `gorm:"column:email"`
	Password    string      `gorm:"column:password"`
	Avatar      *string     `gorm:"column:avatar"`
	Information Information `gorm:"embedded"`
	CreatedAt   int64       `gorm:"column:created_at;autoCreateTime:milli;<-:create"` // allow read and create
	UpdatedAt   int64       `gorm:"column:updated_at;autoCreateTime:milli;autoCreateTime:milli"`
}

func (u *User) TableName() string {
	return "users"
}

// hash password with bcrypt
func (u *User) HashPassword() error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "User.HashPassword.GenerateFromPassword")
	}
	u.Password = string(hashedPass)
	return nil
}

// prepare create user
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}
	return nil
}

// compare password
func (u *User) ComparePassword(hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password)); err != nil {
		return errors.Wrap(err, "User.ComparePassword.CompareHashAndPassword")
	}
	return nil
}

// prepare update user
func (u *User) PrepareUpdate() error {
	if u.Information.FirstName != "" {
		u.Information.FirstName = strings.TrimSpace(u.Information.FirstName)
	}

	if u.Information.LastName != "" {
		u.Information.LastName = strings.TrimSpace(u.Information.LastName)
	}

	if u.Email != "" {
		u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	}

	if u.Password != "" {
		u.Password = strings.TrimSpace(u.Password)
		if err := u.HashPassword(); err != nil {
			return err
		}
	}

	if u.Information.City != nil {
		*u.Information.City = strings.TrimSpace(*u.Information.City)
	}

	if u.Information.PhoneNumber != nil {
		*u.Information.PhoneNumber = strings.TrimSpace(*u.Information.PhoneNumber)
	}

	return nil
}
