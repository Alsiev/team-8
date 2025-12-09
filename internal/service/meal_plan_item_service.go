package service

import (
	"errors"
	"healthy_body/internal/models"
	"healthy_body/internal/repository"
	"log/slog"
)

type MealPlanItemsService interface {
	CreateMealPlanItem(req models.CreateMealPlanItemRequest) (*models.MealPlanItem, error)
	GetAllMealPlanItems() ([]models.MealPlanItem, error)
	UpdateMealPlanItem(id uint, req *models.UpdateMealPlanItemRequest) (*models.MealPlanItem, error)
	GetMealPlanItemById(id uint) (*models.MealPlanItem, error)
	DeleteMealPlanItem(id uint) error
}

type mealPlanItemsService struct {
	mealPlanItems repository.MealPlanItemRepository
	logger        *slog.Logger
}

func NewMealPlanItemsService(
	mealPlanItems repository.MealPlanItemRepository,
	logger *slog.Logger,
) MealPlanItemsService {
	return &mealPlanItemsService{
		mealPlanItems: mealPlanItems,
		logger:        logger,
	}
}

func (s *mealPlanItemsService) CreateMealPlanItem(req models.CreateMealPlanItemRequest) (*models.MealPlanItem, error) {
	if req.MealPlanId == 0 {
		s.logger.Warn("attempt to create item with empty meal plan id")
		return nil, errors.New("meal_plan_id is required")
	}
	if req.Name == "" {
		s.logger.Warn("attempt to create item with empty name")
		return nil, errors.New("name is required")
	}

	item := &models.MealPlanItem{
		Name:        req.Name,
		Description: req.Description,
		Calories:    req.Calories,
		Protein:     req.Protein,
		Carbs:       req.Carbs,
		MealPlanId:  req.MealPlanId,
	}

	if err := s.mealPlanItems.Create(item); err != nil {
		s.logger.Error("failed to create meal plan item", "err", err)
		return nil, err
	}

	s.logger.Info("meal plan item created successfully")
	return item, nil
}

func (s *mealPlanItemsService) GetAllMealPlanItems() ([]models.MealPlanItem, error) {
	mealPlanItems, err := s.mealPlanItems.List()
	if err != nil {
		s.logger.Error("failed to fetch meal plan items", "err", err)
		return nil, err
	}

	if len(mealPlanItems) == 0 {
		s.logger.Warn("service: no meal plans")
	}
	s.logger.Info("meal plan items to fetch successfully", "count", len(mealPlanItems))
	return mealPlanItems, nil
}

func (s *mealPlanItemsService) UpdateMealPlanItem(id uint, req *models.UpdateMealPlanItemRequest) (*models.MealPlanItem, error) {
	mealPlanItems, err := s.mealPlanItems.GetMealPlanItemByID(id)
	if err != nil {
		s.logger.Error("service: meal plan item not found")
		return nil, err
	}

	if req.Name != nil {
		mealPlanItems.Name = *req.Name
	}
	if req.Description != nil {
		mealPlanItems.Description = *req.Description
	}
	if req.Calories != nil {
		mealPlanItems.Calories = *req.Calories
	}
	if req.Protein != nil {
		mealPlanItems.Protein = *req.Protein
	}
	if req.Carbs != nil {
		mealPlanItems.Carbs = *req.Carbs
	}
	if req.MealPlanId != nil {
		mealPlanItems.MealPlanId = *req.MealPlanId
	}

	if err := s.mealPlanItems.Update(mealPlanItems); err != nil {
		s.logger.Error("failed to update meal plan items", "id", id)
		return nil, err
	}
	return mealPlanItems, nil
}

func (s *mealPlanItemsService) GetMealPlanItemById(id uint) (*models.MealPlanItem, error) {
	if id == 0 {
		s.logger.Warn("attempt to meal plan item with id = 0")
		return nil, errors.New("invalid id")
	}
	mealPlanItem, err := s.mealPlanItems.GetMealPlanItemByID(id)
	if err != nil {
		s.logger.Error("failed to get meal plan", "id", id, "error", err)
		return nil, err
	}
	return mealPlanItem, nil
}

func (s *mealPlanItemsService) DeleteMealPlanItem(id uint) error {
	if id == 0 {
		s.logger.Warn("attempt to delete meal plan with id = 0")
		return errors.New("invalid id")
	}
	err := s.mealPlanItems.Delete(id)
	if err != nil {
		s.logger.Error("failed to delete meal plan", "id", id)
		return err
	}
	s.logger.Info("meal plan item deleted successfully", "id", id)
	return nil
}
