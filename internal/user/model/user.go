package model

import (
	"fmt"
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
	Avatar      string      `gorm:"column:avatar"`
	Information Information `gorm:"embedded"`
	CreatedAt   int64       `gorm:"column:created_at;autoCreateTime:milli;<-:create"` // allow read and create
	UpdatedAt   int64       `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
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
	fmt.Println(u.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password)); err != nil {
		return errors.Wrap(err, "User.ComparePassword.CompareHashAndPassword")
	}
	return nil
}

// prepare update user
func (u *User) PrepareUpdate(oldData *User) error {
	if u.Information.FirstName != "" {
		oldData.Information.FirstName = strings.TrimSpace(u.Information.FirstName)
	}

	if u.Information.LastName != "" {
		oldData.Information.LastName = strings.TrimSpace(u.Information.LastName)
	}

	if u.Email != "" {
		oldData.Email = strings.ToLower(strings.TrimSpace(u.Email))
	}

	//if u.Password != "" {
	//	u.Password = strings.TrimSpace(u.Password)
	//	if err := u.HashPassword(); err != nil {
	//		return err
	//	}
	//}

	if u.Information.City != "" {
		oldData.Information.City = strings.TrimSpace(u.Information.City)
	}

	if u.Information.PhoneNumber != "" {
		oldData.Information.PhoneNumber = strings.TrimSpace(u.Information.PhoneNumber)
	}

	return nil
}
