package models

import "time"

type Transaction struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"not null"`
	Amount          float64   `gorm:"not null"`
	TransactionTime time.Time `gorm:"type:timestamp;not null"`
}

type CreditRequest struct {
	Amount float64 `json:"Amount" validate:"required,min=1"`
}
