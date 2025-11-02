package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionExpense TransactionType = "expense"
	TransactionIncome  TransactionType = "income"
)

type Transaction struct {
	ID              int64           `json:"id"`
	UserID          uuid.UUID       `json:"-"`
	WalletID        int64           `json:"wallet_id"`
	CategoryID      int64           `json:"category_id"`
	Amount          int64           `json:"amount"`
	Type            TransactionType `json:"type"`
	Description     *string         `json:"description,omitempty"`
	TransactionDate time.Time       `json:"transaction_date"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// TODO: Tambahkan data join
	// CategoryName string `json:"category_name,omitempty"`
	// WalletName   string `json:"wallet_name,omitempty"`
}

type CreateTransactionRequest struct {
	WalletID        int64      `json:"wallet_id" binding:"required,gt=0"`
	CategoryID      int64      `json:"category_id" binding:"required,gt=0"`
	Amount          int64      `json:"amount" binding:"required,gt=0"`
	Type            string     `json:"type" binding:"required,oneof=expense income"`
	Description     *string    `json:"description"`
	TransactionDate *time.Time `json:"transaction_date"`
}
