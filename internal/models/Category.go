package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type UpdateCategoryRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
}
