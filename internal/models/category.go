package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        int64     `json:"id"`
	UserID    uuid.UUID `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpsertCategoryRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}
