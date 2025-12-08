package repository

import (
	"errors"
	"healthy_body/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type MealPlanRepository interface {
	Create(mealPlan *models.MealPlan) error
	List() ([]models.MealPlan, error)
	Update(mealPlan *models.MealPlan) error
	GetMealPlanByID(id uint) (*models.MealPlan, error)
	Delete(id uint) error
}

type gormMealPlanRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewMealPlanRepository(db *gorm.DB, logger *slog.Logger) MealPlanRepository {
	return &gormMealPlanRepository{
		db:     db,
		logger: logger,
	}
}

func (r *gormMealPlanRepository) Create(mealPlan *models.MealPlan) error {
	if mealPlan == nil {
		r.logger.Warn("attempt to create nil plan")
		return errors.New("plan is nil")
	}
	if err := r.db.Create(mealPlan).Error; err != nil {
		r.logger.Error("failed to create meal plan", "err", err)
		return err
	}
	r.logger.Info("meal plan created")
	return nil
}

func (r *gormMealPlanRepository) List() ([]models.MealPlan, error) {
	var mealPlans []models.MealPlan

	if err := r.db.Preload("Meals").Find(&mealPlans).Error; err != nil {
		r.logger.Error("failed to fetch meal plans")
		return nil, err
	}
	r.logger.Info("meal plans fetched", "count", len(mealPlans))
	return mealPlans, nil
}

func (r *gormMealPlanRepository) Update(mealPlan *models.MealPlan) error {
	if mealPlan == nil {
		r.logger.Warn("attempt to update nil meal plan")
		return errors.New("meal plan is nil")
	}
	err := r.db.Model(&models.MealPlan{}).Where("id =?", mealPlan.ID).
		Select("CategoryID", "TotalDays").Updates(mealPlan).Error

	if err != nil {
		r.logger.Error("failed to update meal plan", "id", mealPlan.ID, "err", err)
		return err
	}
	r.logger.Info("meal plan updated successfully", "id", mealPlan.ID)
	return nil
}

func (r *gormMealPlanRepository) GetMealPlanByID(id uint) (*models.MealPlan, error) {
	var mealPlan models.MealPlan

	if err := r.db.Preload("Meals").First(&mealPlan, id).Error; err != nil {
		r.logger.Error("failed to fetch meal plan", "err", err)
		return nil, err
	}
	r.logger.Info("fetch to meal plan successfully", "id", id)
	return &mealPlan, nil
}

func (r *gormMealPlanRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.MealPlan{}, id).Error; err != nil {
		r.logger.Error("failed to delete meal plan", "err", err)
		return err
	}
	r.logger.Info("meal plan delete successfully")
	return nil
}
