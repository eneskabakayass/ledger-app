package models

import (
	"ledger-app/internal/validation"
)

type User struct {
	ID           uint          `gorm:"primaryKey"`
	Name         string        `gorm:"not null" validate:"required,min=1,max=10"`
	PasswordHash string        `gorm:"not null"`
	IsAdmin      bool          `gorm:"default:false"`
	Credits      []Transaction `gorm:"foreignKey:UserID"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct().Struct(u)
}
