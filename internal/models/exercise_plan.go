package models

import "gorm.io/gorm"

type ExercisePlan struct {
	gorm.Model
	Name          string `json:"name"`
	Description   string `json:"description"`
	DurationWeeks int    `json:"duration_weeks"`

	Exercises []ExercisePlanItem `json:"exercises" gorm:"foreignKey:ExercisePlanID"`

	CategoriesID uint        `json:"categories_id"`
	Categories   *Categories `json:"-"`
}

type CreateExercesicePlanRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	CategoryID    uint `json:"categories_id"`
	DurationWeeks int  `json:"duration_weeks"`
}

type UpdateExercesicePlanRequest struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	DurationWeeks *int `json:"duration_weeks"`
}
