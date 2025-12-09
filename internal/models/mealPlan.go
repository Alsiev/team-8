package models

import "gorm.io/gorm"

type MealPlan struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CategoryID  uint           `json:"category_id"`
	TotalDays   int            `json:"total_days"`
	Meals       []MealPlanItem `json:"meals"`
	Category    *Category      `json:"-"`
}

type CreateMealPlanRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  uint   `json:"category_id"`
	TotalDays   int    `json:"total_days"`
}

type UpdateMealPlanRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	CategoryID  *uint   `json:"category_id"`
	TotalDays   *int    `json:"total_days"`
}
