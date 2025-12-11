package service

import (
	"errors"
	"healthy_body/internal/models"
	"healthy_body/internal/repository"
	"log/slog"
)

type MealPlanService interface {
	CreateMealPlan(req models.CreateMealPlanRequest) (*models.MealPlan, error)
	ListMealPlan() ([]models.MealPlan, error)
	UpdateMealPlan(id uint, req *models.UpdateMealPlanRequest) (*models.MealPlan, error)
	GetMealPlanByID(id uint) (*models.MealPlan, error)
	DeleteMealPlan(id uint) error
}

type mealPlanService struct {
	mealPlans repository.MealPlanRepository
	logger    *slog.Logger
	category  CategoryServices
}

func NewMealPlanService(
	mealPlans repository.MealPlanRepository,
	logger *slog.Logger,
	category CategoryServices,
) MealPlanService {
	return &mealPlanService{
		mealPlans: mealPlans,
		logger:    logger,
		category:  category,
	}
}

func (s *mealPlanService) CreateMealPlan(req models.CreateMealPlanRequest) (*models.MealPlan, error) {
	if *req.CategoriesID == 0 {
		s.logger.Error("invalid category id", "id", req.CategoriesID)
		return nil, errors.New("category id is required")
	}
	if req.TotalDays <= 0 {
		s.logger.Error("invalid total_days")
		return nil, errors.New("total days must be greater than zero")
	}

	if _, err := s.category.GetCategoryByID(*req.CategoriesID); err != nil {
		s.logger.Error("error GetCategoryByID function in exercise_service.go")
		return nil, err
	}

	mealPlan := models.MealPlan{
		Name:        req.Name,
		Description: req.Description,
		CategoriesID:  req.CategoriesID,
		TotalDays:   req.TotalDays,
	}

	if err := s.mealPlans.Create(&mealPlan); err != nil {
		s.logger.Error("service: failed to create meal plan")
		return nil, err
	}
	return &mealPlan, nil
}

func (s *mealPlanService) ListMealPlan() ([]models.MealPlan, error) {
	mealPlans, err := s.mealPlans.List()
	if err != nil {
		s.logger.Error("failed to fetch meal plans")
		return nil, err
	}

	if len(mealPlans) == 0 {
		s.logger.Warn("service: no meal plans")
	}

	s.logger.Info("meal plans fetch to successfully", "count", len(mealPlans))
	return mealPlans, nil
}

func (s *mealPlanService) UpdateMealPlan(id uint, req *models.UpdateMealPlanRequest) (*models.MealPlan, error) {
	mealPlan, err := s.mealPlans.GetMealPlanByID(id)
	if err != nil {
		s.logger.Error("service: meal plan not found")
		return nil, err
	}

	if req.Name != nil {
		mealPlan.Name = *req.Name
	}

	if req.Description != nil {
		mealPlan.Description = *req.Description
	}

	if req.CategoriesID != nil {
		mealPlan.CategoriesID = req.CategoriesID
	}
	if req.TotalDays != nil {
		mealPlan.TotalDays = *req.TotalDays
	}

	if err := s.mealPlans.Update(mealPlan); err != nil {
		s.logger.Error("failed to update meal plan", "id", id)
		return nil, err
	}
	return mealPlan, nil
}

func (s *mealPlanService) DeleteMealPlan(id uint) error {
	if id == 0 {
		s.logger.Warn("attempt to delete meal plan with id = 0")
		return errors.New("invalid id")
	}
	err := s.mealPlans.Delete(id)
	if err != nil {
		s.logger.Error("failed to delete meal plan", "id", id)
		return err
	}
	s.logger.Info("meal plan deleted successfully", "id", id)
	return nil
}

func (s *mealPlanService) GetMealPlanByID(id uint) (*models.MealPlan, error) {
	if id == 0 {
		s.logger.Warn("attempt to meal plan with id = 0")
		return nil, errors.New("invalid id")
	}
	mealPlan, err := s.mealPlans.GetMealPlanByID(id)
	if err != nil {
		s.logger.Error("failed to get meal plan", "id", id, "error", err)
		return nil, err
	}
	return mealPlan, nil
}
