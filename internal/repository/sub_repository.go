package repository

import (
	"errors"
	"healthy_body/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type SubscriptionRepo interface {
	Create(req *models.Subscription) error
	GetByID(id uint) (*models.Subscription, error)
	GetList() ([]models.Subscription, error)
	Update(up *models.Subscription) error
	Delete(id uint) error
}

type subscriptionRepo struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewSubscriptionRepo(db *gorm.DB, log *slog.Logger) SubscriptionRepo {
	return &subscriptionRepo{
		db:  db,
		log: log,
	}
}

func (r *subscriptionRepo) Create(req *models.Subscription) error {
	if req == nil {
		r.log.Error("error create function in sub_repository.go")
		return errors.New("error create sub in db")
	}

	return r.db.Create(req).Error
}

func (r *subscriptionRepo) GetByID(id uint) (*models.Subscription, error) {
	var sub models.Subscription
	if err := r.db.First(&sub, id).Error; err != nil {
		r.log.Error("error getbyid function in sub_repository.go")
		return nil, err
	}

	return &sub, nil
}

func (r *subscriptionRepo) GetList() ([]models.Subscription, error) {
	var sub []models.Subscription
	if err := r.db.Find(&sub).Error; err != nil {
		r.log.Error("error getList function in sub_repository.go")
		return nil, err
	}

	return sub, nil
}

func (r *subscriptionRepo) Update(up *models.Subscription) error {
	if up == nil {
		r.log.Error("error update function in sub_repository.go")
		return errors.New("not found sub to update in db")
	}

	return r.db.Save(up).Error
}

func (r *subscriptionRepo) Delete(id uint) error {
	if err := r.db.Delete(&models.Subscription{}, id).Error; err != nil {
		r.log.Error("error delete function in sub_repository.go")
		return errors.New("error delete sub by id in dn")
	}

	return nil
}
