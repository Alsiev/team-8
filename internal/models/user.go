package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`

	CategoryID uint      `json:"category_id"`
	Category   *Category `json:"-"`
}

type CreateUserRequest struct {
	Name    string  `json:"name"`
}

type UpdateUserRequest struct {
	Name    *string  `json:"name"`
	Balance *float64 `json:"balance"`
}
