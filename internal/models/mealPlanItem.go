package models


type MealPlanItem struct {
	ID         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string     `json:"name"`
	Calories   float64    `json:"calories"`
	Protein    float64    `json:"protein"`
	Carbs      float64    `json:"carbs"`

	MealPlanId uint       `json:"meal_plan_id"`
	MealPlan   *MealPlan  `json:"-"`
}

type CreateMealPlanItemRequest struct {
	Name       string  `json:"name"`
	Calories   float64 `json:"calories"`
	Protein    float64 `json:"protein"`
	Carbs      float64 `json:"carbs"`
	MealPlanId uint    `json:"meal_plan_id"`
}

type UpdateMealPlanItemRequest struct {
	Name       string  `json:"name"`
	Calories   float64 `json:"calories"`
	Protein    float64 `json:"protein"`
	Carbs      float64 `json:"carbs"`
	MealPlanId uint    `json:"meal_plan_id"`
}

