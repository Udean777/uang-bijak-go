package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        int64     `json:"id"`
	UserID    uuid.UUID `json:"-"`
	Name      string    `json:"name"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateWalletRequest struct {
	Name           string `json:"name" binding:"required,min=3,max=100"`
	InitialBalance int64  `json:"initial_balance" binding:"gte=0"`
}

type UpdateWalletRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}
