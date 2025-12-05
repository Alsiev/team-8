package models

type ExercisePlanItem struct {
	ID              uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string `json:"name"`
	Sets            int    `json:"sets"`
	Reps            int    `json:"reps"`
	DurationMinutes string `json:"duration_minutes"`
	EquipmentNeeded string `json:"equipment_needed"`
	DayOfWeek       string `json:"day_of_week"`

	ExercisePlanID uint          `json:"exercise_plan_id"`
	ExercisePlan   *ExercisePlan `json:"-"`
}

type CreateExercisePlanItemRequest struct {
	Name            string `json:"name"`
	Sets            int    `json:"sets"`
	Reps            int    `json:"reps"`
	DurationMinutes string `json:"duration_minutes"`
	EquipmentNeeded string `json:"equipment_needed"`
	DayOfWeek       string `json:"day_of_week"`
	ExercisePlanID  uint   `json:"exercise_plan_id"`
}

type UpdateExercisePlanItemRequest struct {
	Name            *string `json:"name"`
	Sets            *int    `json:"sets"`
	Reps            *int    `json:"reps"`
	DurationMinutes *string `json:"duration_minutes"`
	EquipmentNeeded *string `json:"equipment_needed"`
	DayOfWeek       *string `json:"day_of_week"`
	ExercisePlanID  *uint   `json:"exercise_plan_id"`
}
