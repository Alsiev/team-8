package models

import "gorm.io/gorm"

type Subscription struct {
	gorm.Model
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Price        int         `json:"price"`
	DurationDays int         `json:"duration_days"`
	CategoriesID uint        `json:"categories_id"`
	Categories   *Categories `json:"-"`
}

type CreateSubscriptionRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Price        int    `json:"price"`
	DurationDays int    `json:"duration_days"`
	CategoriesID uint   `json:"categories_id"`
}

type UpdateSubscriptionRequest struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	Price        *int    `json:"price"`
	DurationDays *int    `json:"duration_days"`
	CategoriesID *uint   `json:"categories_id"`
}
