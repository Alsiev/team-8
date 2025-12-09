package models

import "gorm.io/gorm"

type MealPlanItem struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Calories    float64   `json:"calories"`
	Protein     float64   `json:"protein"`
	Carbs       float64   `json:"carbs"`
	MealPlanId  uint      `json:"meal_plan_id"`
	MealPlan    *MealPlan `json:"-"`
}

type CreateMealPlanItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Calories    float64 `json:"calories"`
	Protein     float64 `json:"protein"`
	Carbs       float64 `json:"carbs"`
	MealPlanId  uint    `json:"meal_plan_id"`
}

type UpdateMealPlanItemRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Calories    *float64 `json:"calories"`
	Protein     *float64 `json:"protein"`
	Carbs       *float64 `json:"carbs"`
	MealPlanId  *uint    `json:"meal_plan_id"`
}
