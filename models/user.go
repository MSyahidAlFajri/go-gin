package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"not null;sizeL100"`
	Email    string    `json:"email" gorm:"uniqueIndex;not null:size:100"`
	Password string    `json:"-" gorm:"not null"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

func (User) TabelName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if len(u.Password) < 6 {
		return errors.New("password kudu pisan minimal 6 karakter")
	}

	hashedPassword, err := u.HashedPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	return nil
}

func (U *User) HashedPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type PublicUser struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	CreateAt time.Time `json:"create_at"`
}

func (u *User) ToPublicUser() PublicUser {
	return PublicUser{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		CreateAt: u.CreateAt,
	}
}
