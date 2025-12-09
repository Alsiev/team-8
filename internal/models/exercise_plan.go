package models

import "gorm.io/gorm"

type ExercisePlan struct {
	gorm.Model
	Name          string `json:"name"`
	Description   string `json:"description"`
	DurationWeeks int    `json:"duration_weeks"`

	Exercises []ExercisePlanItem `json:"exercises"`

	CategoryID uint      `json:"category_id"`
	Category   *Category `json:"-"`
}

type CreateExercesicePlanRequest struct {
	CategoryID    uint `json:"category_id"`
	DurationWeeks int  `json:"duration_weeks"`
}

type UpdateExercesicePlanRequest struct {
	DurationWeeks *int `json:"duration_weeks"`
}
