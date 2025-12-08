package repository

import (
	"errors"
	"healthy_body/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type ExercisePlanRepo interface {
	CreateExercisePlan(exercise *models.ExercisePlan) error
	GetByIDExercisePlan(id uint) (*models.ExercisePlan, error)
	GetByIDExercisePlanForNotPreload(id uint) (*models.ExercisePlan, error)
	GetAllExercisePlan() ([]models.ExercisePlan, error)
	UpdateExercisePlan(exercise *models.ExercisePlan) error
	DeleteExercisePlan(id uint) error

	CreateExercisePlanItem(item *models.ExercisePlanItem) error
	GetAllExercisePlanItem() ([]models.ExercisePlanItem, error)
	GetByIDExercisePlanItem(id uint) (*models.ExercisePlanItem, error)
	UpdateExercisePlanItem(exercise *models.ExercisePlanItem) error
	DeleteExercisePlanItem(id uint) error
}

type exercisePlanRepo struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewExercisePlanRepo(db *gorm.DB, log *slog.Logger) ExercisePlanRepo {
	return &exercisePlanRepo{
		db:  db,
		log: log,
	}
}

func (r *exercisePlanRepo) CreateExercisePlan(exercise *models.ExercisePlan) error {
	if exercise == nil {
		r.log.Error("error in Create function exercise_plan_repository.go")
		return errors.New("error create in db")
	}

	return r.db.Create(exercise).Error
}

func (r *exercisePlanRepo) GetByIDExercisePlan(id uint) (*models.ExercisePlan, error) {
	var exercise models.ExercisePlan

	if err := r.db.Preload("Exercises").Preload("Category").First(&exercise, id).Error; err != nil {
		r.log.Error("error in GetByID function exercise_plan_repository.go")
		return nil, err
	}

	return &exercise, nil
}

func (r *exercisePlanRepo) GetByIDExercisePlanForNotPreload(id uint) (*models.ExercisePlan, error) {
	var exercise models.ExercisePlan

	if err := r.db.First(&exercise, id).Error; err != nil {
		r.log.Error("error in GetByID function exercise_plan_repository.go")
		return nil, err
	}

	return &exercise, nil
}

func (r *exercisePlanRepo) GetAllExercisePlan() ([]models.ExercisePlan, error) {
	var exercises []models.ExercisePlan
	if err := r.db.Find(&exercises).Error; err != nil {
		r.log.Error("error in GetAll function exercise_plan_repository.go")
		return nil, err
	}
	return exercises, nil
}

func (r *exercisePlanRepo) UpdateExercisePlan(exercise *models.ExercisePlan) error {
	if exercise == nil {
		r.log.Error("error in Update function exercise_plan_repository.go")
		return errors.New("error update in db")
	}

	return r.db.Save(exercise).Error
}

func (r *exercisePlanRepo) DeleteExercisePlan(id uint) error {
	if err := r.db.Delete(&models.ExercisePlan{}, id).Error; err != nil {
		r.log.Error("error in Delete function exercise_plan_repository.go")
		return errors.New("error delete in db")
	}

	return nil
}

func (r *exercisePlanRepo) CreateExercisePlanItem(item *models.ExercisePlanItem) error {
	if item == nil {
		r.log.Error("error in Create function exercise_plan_item_repository.go")
		return errors.New("error create in db")
	}

	return r.db.Create(item).Error
}

func (r *exercisePlanRepo) GetAllExercisePlanItem() ([]models.ExercisePlanItem, error) {
	var exercises []models.ExercisePlanItem
	if err := r.db.Find(&exercises).Error; err != nil {
		r.log.Error("error in GetAll function exercise_plan_item_repository.go")
		return nil, err
	}
	return exercises, nil
}

func (r *exercisePlanRepo) GetByIDExercisePlanItem(id uint) (*models.ExercisePlanItem, error) {
	var exercise models.ExercisePlanItem

	if err := r.db.First(&exercise, id).Error; err != nil {
		r.log.Error("error in GetByID function exercise_plan_repository.go")
		return nil, err
	}

	return &exercise, nil
}

func (r *exercisePlanRepo) UpdateExercisePlanItem(exercise *models.ExercisePlanItem) error {
	if exercise == nil {
		r.log.Error("error in Update function exercise_plan_item_repository.go")
		return errors.New("error update in db")
	}

	return r.db.Save(exercise).Error
}

func (r *exercisePlanRepo) DeleteExercisePlanItem(id uint) error {
	if err := r.db.Delete(&models.ExercisePlanItem{}, id).Error; err != nil {
		r.log.Error("error in Delete function exercise_plan_item_repository.go")
		return errors.New("error delete in db")
	}

	return nil
}
