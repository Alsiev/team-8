package service

import (
	"errors"
	"healthy_body/internal/models"
	"healthy_body/internal/repository"
	"log/slog"
)

type SubscriptionService interface {
	CreateSub(req *models.CreateSubscriptionRequest) (*models.Subscription, error)
	GetSubByID(id uint) (*models.Subscription, error)
	GetListSub() ([]models.Subscription, error)
	UpdateSub(id uint, req models.UpdateSubscriptionRequest) (*models.Subscription, error)
	Delete(id uint) error
}

type subscriptionService struct {
	subRepo  repository.SubscriptionRepo
	category CategoryServices
	log      *slog.Logger
}

func NewSubscriptionService(subRepo repository.SubscriptionRepo, log *slog.Logger, category CategoryServices) SubscriptionService {
	return &subscriptionService{
		subRepo:  subRepo,
		log:      log,
		category: category,
	}
}

func (s *subscriptionService) CreateSub(req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	if err := s.valiD(req); err != nil {
		s.log.Error("error valid sub struct")
		return nil, err
	}

	if _, err := s.category.GetCategoryByID(req.CategoryID); err != nil {
		s.log.Error("error found category id")
		return nil, err
	}

	sub := &models.Subscription{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		CategoryID: req.CategoryID,
		DurationDays: req.DurationDays,
	}

	if err := s.subRepo.Create(sub); err != nil {
		s.log.Error("error create sub in sub_service.go")
		return nil, err
	}

	return sub, nil
}

func (s *subscriptionService) GetSubByID(id uint) (*models.Subscription, error) {
	sub, err := s.subRepo.GetByID(id)
	if err != nil {
		s.log.Error("error not found id")
		return nil, err
	}

	return sub, err
}

func (s *subscriptionService) GetListSub() ([]models.Subscription, error) {
	list, err := s.subRepo.GetList()
	if err != nil {
		s.log.Error("")
		return nil, err
	}

	return list, err
}

func (s *subscriptionService) UpdateSub(id uint, req models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	sub, err := s.subRepo.GetByID(id)
	if err != nil {
		s.log.Error("error GetByID function")
		return nil, err
	}

	s.upSub(sub, req)
	if err := s.subRepo.Update(sub); err != nil {
		s.log.Error("error update function in sub_service.go")
		return nil, err
	}

	return sub, nil
}

func (s *subscriptionService) Delete(id uint) error {
	if err := s.subRepo.Delete(id); err != nil {
		s.log.Error("error not found sub by id for delete or delete error")
		return err
	}

	return nil
}

func (s *subscriptionService) valiD(req *models.CreateSubscriptionRequest) error {
	if req.Name == "" {
		return errors.New("empty name in sub struct")
	}
	if req.Description == "" {
		return errors.New("empty description in sub struct")
	}

	if req.Price == 0 {
		return errors.New("empty price in sub struct")
	}
	if req.DurationDays == 0 {
		return errors.New("empty duration days in sub struct")
	}

	return nil
}

func (s *subscriptionService) upSub(sub *models.Subscription, req models.UpdateSubscriptionRequest) {
	if req.Name != nil {
		sub.Name = *req.Name
	}

	if req.Description != nil {
		sub.Description = *req.Description
	}

	if req.Price != nil {
		sub.Price = *req.Price
	}

	if req.DurationDays != nil {
		sub.DurationDays = *req.DurationDays
	}
}
