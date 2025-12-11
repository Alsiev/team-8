package service

import (
	"errors"
	"healthy_body/internal/models"
	"healthy_body/internal/repository"
	"log/slog"
)

type ExercisePlanServices interface {
	CreatePlan(req models.CreateExercesicePlanRequest) (*models.ExercisePlan, error)
	GetPlanByID(id uint) (*models.ExercisePlan, error)
	GetPlanByIDNotPreloads(id uint) (*models.ExercisePlan, error)
	GetListPlans() ([]models.ExercisePlan, error)
	UpdatePlan(id uint, req models.UpdateExercesicePlanRequest) (*models.ExercisePlan, error)
	DeletePlan(id uint) error

	CreatePlanItem(req models.CreateExercisePlanItemRequest) (*models.ExercisePlanItem, error)
	GetAllPlanItem() ([]models.ExercisePlanItem, error)
	GetByIDPlanItem(id uint) (*models.ExercisePlanItem, error)
	UpdatePlanItem(id uint, req models.UpdateExercisePlanItemRequest) (*models.ExercisePlanItem, error)
	DeletePlanItem(id uint) error
}

type exercisePlanServices struct {
	exerciseRepo repository.ExercisePlanRepo
	category     CategoryServices
	log          *slog.Logger
}

func NewExercisePlanServices(exerciseRepo repository.ExercisePlanRepo, log *slog.Logger, category CategoryServices) ExercisePlanServices {
	return &exercisePlanServices{
		exerciseRepo: exerciseRepo,
		log:          log,
		category: category,
	}
}

func (e *exercisePlanServices) CreatePlan(req models.CreateExercesicePlanRequest) (*models.ExercisePlan, error) {
	if req.DurationWeeks == 0 {
		e.log.Error("error CreatePlan function in exercise_service.go")
		return nil, errors.New("empty weeks your plan")
	}

	if _, err := e.category.GetCategoryByID(req.CategoryID); err != nil {
		e.log.Error("error GetCategoryByID function in exercise_service.go")
		return nil, err
	}

	exercise := &models.ExercisePlan{
		Name:  req.Name,
		Description: req.Description,
		DurationWeeks: req.DurationWeeks,
		CategoriesID:    req.CategoryID,
	}

	if err := e.exerciseRepo.CreateExercisePlan(exercise); err != nil {
		e.log.Error("error CreatePlan function in exercise_service.go")
		return nil, err
	}

	return exercise, nil
}

func (e *exercisePlanServices) GetPlanByID(id uint) (*models.ExercisePlan, error) {
	plan, err := e.exerciseRepo.GetByIDExercisePlan(id)
	if err != nil {
		e.log.Error("error GetPlanByID function in exercise_service.go")
		return nil, err
	}

	return plan, nil
}

func (e *exercisePlanServices) GetPlanByIDNotPreloads(id uint) (*models.ExercisePlan, error) {
	plan, err := e.exerciseRepo.GetByIDExercisePlan(id)
	if err != nil {
		e.log.Error("error GetPlanByID function in exercise_service.go")
		return nil, err
	}

	return plan, nil
}

func (e *exercisePlanServices) GetListPlans() ([]models.ExercisePlan, error) {
	list, err := e.exerciseRepo.GetAllExercisePlan()
	if err != nil {
		e.log.Error("error GetListPlans function in exercise_service.go")
		return nil, err
	}

	return list, nil
}

func (e *exercisePlanServices) UpdatePlan(id uint, req models.UpdateExercesicePlanRequest) (*models.ExercisePlan, error) {
	plan, err := e.GetPlanByID(id)
	if err != nil {
		e.log.Error("error UpdatePlan function in exercise_service.go")
		return nil, err
	}

	if req.DurationWeeks != nil {
		plan.DurationWeeks = *req.DurationWeeks
	}

	if req.Name != nil {
		plan.Name = *req.Name
	}

	if req.Description != nil {
		plan.Description = *req.Description
	}

	if err := e.exerciseRepo.UpdateExercisePlan(plan); err != nil {
		e.log.Error("error UpdatePlan function in exercise_service.go")
		return nil, err
	}

	return plan, nil
}

func (e *exercisePlanServices) DeletePlan(id uint) error {
	if err := e.exerciseRepo.DeleteExercisePlan(id); err != nil {
		e.log.Error("error DeletePlan function in exercise_service.go")
		return err
	}

	return nil
}


func (e *exercisePlanServices) CreatePlanItem(req models.CreateExercisePlanItemRequest) (*models.ExercisePlanItem, error) {
	if err := e.validate(req); err != nil {
		e.log.Error("error CreatePlanItem function in exercise_service.go")
		return nil, err
	}

	if _, err := e.exerciseRepo.GetByIDExercisePlanForNotPreload(req.ExercisePlanID); err != nil {
		e.log.Error("error CreatePlanItem function in exercise_service.go")
		return nil, err
	}

	item := &models.ExercisePlanItem{
		Name:            req.Name,
		Sets:            req.Sets,
		Reps:            req.Reps,
		EquipmentNeeded: req.EquipmentNeeded,
		DurationMinutes: req.DurationMinutes,
		DayOfWeek:       req.DayOfWeek,
		ExercisePlanID:  req.ExercisePlanID,
	}

	if err := e.exerciseRepo.CreateExercisePlanItem(item); err != nil {
		e.log.Error("error CreatePlanItem function in exercise_service.go")
		return nil, err
	}

	return item, nil
}

func (e *exercisePlanServices) GetAllPlanItem() ([]models.ExercisePlanItem, error) {
	item, err := e.exerciseRepo.GetAllExercisePlanItem()
	if err != nil {
		e.log.Error("error GetAllPlanItem function in exercise_service.go")
		return nil, err
	}

	return item, nil
}

func (e *exercisePlanServices) GetByIDPlanItem(id uint) (*models.ExercisePlanItem, error) {
	item, err := e.exerciseRepo.GetByIDExercisePlanItem(id)
	if err != nil {
		e.log.Error("error GetByIDPlanItem function in exercise_service.go")
		return nil, err
	}

	return item, nil
}

func (e *exercisePlanServices) UpdatePlanItem(id uint, req models.UpdateExercisePlanItemRequest) (*models.ExercisePlanItem, error) {
	item, err := e.exerciseRepo.GetByIDExercisePlanItem(id)
	if err != nil {
		e.log.Error("error UpdatePlanItem function in exercise_service.go")
		return nil, err
	}

	e.up(item, req)

	if err := e.exerciseRepo.UpdateExercisePlanItem(item); err != nil {
		e.log.Error("error UpdatePlanItem function in exercise_service.go")
		return nil, err
	}

	return item, nil
}

func (e *exercisePlanServices) DeletePlanItem(id uint) error {
	if err := e.exerciseRepo.DeleteExercisePlanItem(id); err != nil {
		e.log.Error("error DeletePlanItem function in exercise_service.go")
		return err
	}

	return nil
}

func (r *exercisePlanServices) validate(req models.CreateExercisePlanItemRequest) error {
	if req.Name == "" {
		return errors.New("name plan item is null")
	}
	if req.Sets == 0 {
		return errors.New("sets plan item is null")
	}

	if req.Reps == 0 {
		return errors.New("reps plan item is null")
	}

	if req.DurationMinutes == "" {
		return errors.New("durationMinutes plan item is null")
	}

	if req.DayOfWeek == "" {
		return errors.New("dayOfWeek plan item is null")
	}

	if req.EquipmentNeeded == "" {
		return errors.New("equipmentNeeded plan item is null")
	}

	return nil
}

func (r *exercisePlanServices) up(item *models.ExercisePlanItem, req models.UpdateExercisePlanItemRequest) {
	if req.Name != nil {
		item.Name = *req.Name
	}

	if req.Sets != nil {
		item.Sets = *req.Sets
	}

	if req.Reps != nil {
		item.Reps = *req.Reps
	}

	if req.DurationMinutes != nil {
		item.DurationMinutes = *req.DurationMinutes
	}

	if req.EquipmentNeeded != nil {
		item.EquipmentNeeded = *req.EquipmentNeeded
	}

	if req.DayOfWeek != nil {
		item.DayOfWeek = *req.DayOfWeek
	}

}
