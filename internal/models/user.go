package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string    `json:"name"`
	Balanse   float64   `json:"balanse"`

	CategoryID uint     `json:"category_id"`
	Category  *Category `json:"-"`
}

type CreateUserRequest struct {
	Name    string  `json:"name"`
	Balanse float64 `json:"balanse"`
	// CategoryID если нужно — скажи, добавлю
}

type UpdateUserRequest struct {
	Name    *string  `json:"name"`
	Balanse *float64 `json:"balanse"`
}
