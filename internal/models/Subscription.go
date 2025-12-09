package models

import "gorm.io/gorm"

type Subscription struct {
	gorm.Model
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         int       `json:"price"`
	Duration_days int       `json:"duration_days"`
	Category_id   uint      `json:"category_id"`
	Category      *Category `json:"-"`
}

type CreateSubscriptionRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Price         int    `json:"price"`
	Duration_days int    `json:"duration_days"`
	Category_id   uint   `json:"category_id"`
}

type UpdateSubscriptionRequest struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Price         *int    `json:"price"`
	Duration_days *int    `json:"duration_days"`
	Category_id   *uint   `json:"category_id"`
}
