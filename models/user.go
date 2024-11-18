package models

import "github.com/go-playground/validator/v10"

type User struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"not null" validate:"required,min=1,max=10"`
	Credit int    `gorm:"default:0"`
}

type CreditRequest struct {
	Credit int `json:"credit" validate:"required,min=1"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (u *User) Validate() error {
	return validate.Struct(u)
}
