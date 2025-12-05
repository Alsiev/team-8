package models

import "gorm.io/gorm"

type MealPlan struct {
	gorm.Model
	CategoryID uint           `json:"category_id"`
	TotalDays  int            `json:"total_days"`
	Meals      []MealPlanItem `json:"meals"`
}

type CreateMealPlanRequest struct {
	CategoryID uint `json:"category_id"`
	TotalDays  int  `json:"total_days"`
}

type UpdateMealPlanRequest struct {
	CategoryID *uint `json:"category_id"`
	TotalDays  *int  `json:"total_days"`
}
