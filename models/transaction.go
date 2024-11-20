package models

import (
	"ledger-app/internal/validation"
	"time"
)

type Transaction struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"not null"`
	Amount          float64   `gorm:"not null"`
	TransactionTime time.Time `gorm:"type:timestamp;" validate:"required"`
}

type CreditRequest struct {
	Amount float64 `json:"Amount" validate:"required,gt=0"`
}

func (t *Transaction) Validator() error {
	return validation.ValidateStruct().Struct(t)
}
