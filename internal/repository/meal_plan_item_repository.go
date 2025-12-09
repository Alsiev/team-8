package repository

import (
	"errors"
	"healthy_body/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type MealPlanItemRepository interface {
	Create(mealPlanItem *models.MealPlanItem) error
	List() ([]models.MealPlanItem, error)
	Update(mealPlan *models.MealPlanItem) error
	GetMealPlanItemByID(id uint) (*models.MealPlanItem, error)
	Delete(id uint) error
}

type gormMealPlanItemRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewMealPlanItemRepository(db *gorm.DB, logger *slog.Logger) MealPlanItemRepository {
	return &gormMealPlanItemRepository{
		db:     db,
		logger: logger,
	}
}

func (r *gormMealPlanItemRepository) Create(mealPlanItem *models.MealPlanItem) error {
	if mealPlanItem == nil {
		r.logger.Error("failed to create meal plan item")
		return errors.New("meal plan item is nil")
	}
	if err := r.db.Create(mealPlanItem).Error; err != nil {
		r.logger.Error("failed to create meal plan item", "err", err)
		return err
	}

	r.logger.Info("meal plan item created")
	return nil
}

func (r *gormMealPlanItemRepository) List() ([]models.MealPlanItem, error) {
	var mealPlanItems []models.MealPlanItem

	if err := r.db.Find(&mealPlanItems).Error; err != nil {
		r.logger.Error("failed to fetch meal plan items", "err", err)
		return nil, err
	}

	r.logger.Info("fetched meal plan item successfully", "count", len(mealPlanItems))
	return mealPlanItems, nil
}

func (r *gormMealPlanItemRepository) Update(mealPlanItem *models.MealPlanItem) error {
	if mealPlanItem == nil {
		r.logger.Warn("attempt to update nil meal plan item")
		return errors.New("meal plan item is nil")
	}

	err := r.db.Model(&models.MealPlanItem{}).Where("id = ?", mealPlanItem.ID).
		Select("Name", "Calories", "Protein", "Carbs", "MealPlanId").Updates(mealPlanItem).Error

	if err != nil {
		r.logger.Error("failed to update meal plan", "id", mealPlanItem.ID, "err", err)
		return err
	}
	r.logger.Info("meal plan item updated successfully", "id", mealPlanItem.ID)
	return nil
}

func (r *gormMealPlanItemRepository) GetMealPlanItemByID(id uint) (*models.MealPlanItem, error) {
	var mealPlanItem models.MealPlanItem

	if err := r.db.First(&mealPlanItem, id).Error; err != nil {
		r.logger.Error("failed to fetch meal plan item", "err", err)
		return nil, err
	}
	r.logger.Info("fetch meal plan item successfully")
	return &mealPlanItem, nil
}

func (r *gormMealPlanItemRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.MealPlanItem{}, id).Error; err != nil {
		r.logger.Error("failed to delete meal plan item", "id", id)
		return err
	}
	r.logger.Info("meal plan item deleted", "id", id)
	return nil
}
