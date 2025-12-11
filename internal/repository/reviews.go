package repository

import (
	"fmt"
	"healthy_body/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type ReviewsRepository interface {
	CreateReviews(req *models.Reviews) error
	GetAllReviews() ([]models.Reviews, error)
	GetReviewsByID(id uint) (*models.Reviews, error)
	UpdateReviews(req *models.Reviews) error
	Delete(id uint) error

	GetByUserID(userID uint) ([]models.Reviews, error)
	GetByCategoryID(categoryID uint) ([]models.Reviews, error)
}

type reviewsRepository struct {
	reviews *gorm.DB
	log     *slog.Logger
}

func NewReviewsRepository(reviews *gorm.DB, log *slog.Logger) ReviewsRepository {
	return &reviewsRepository{reviews: reviews, log: log}
}

func (r *reviewsRepository) CreateReviews(req *models.Reviews) error {

	if err := r.reviews.Create(req).Error; err != nil {
		r.log.Error("Ошибка создания отзыва",
			"error", err.Error())
		return fmt.Errorf("ошибка создания отзыва: %w", err)
	}

	r.log.Info("Отзыв успешно создан")

	return nil
}

func (r *reviewsRepository) GetAllReviews() ([]models.Reviews, error) {
	var reviews []models.Reviews

	if err := r.reviews.Find(&reviews).Error; err != nil {
		r.log.Error("Ошибка при выдаче всех отзывов",
			"error", err.Error())
		return nil, fmt.Errorf("ошибка при выдаче всех отзывов: %w", err)
	}

	r.log.Info("Отзывы получены успешно")

	return reviews, nil
}

func (r *reviewsRepository) GetReviewsByID(id uint) (*models.Reviews, error) {
	var reviews models.Reviews

	if err := r.reviews.First(&reviews, id).Error; err != nil {
		r.log.Error("Ошибка при выводе отзыва",
			"error", err.Error())
		return nil, fmt.Errorf("ошибка при выводе отзыва %w", err)
	}

	r.log.Info("Отзыв получен")
	return &reviews, nil
}

func (r *reviewsRepository) UpdateReviews(req *models.Reviews) error {
	if err := r.reviews.Model(&req).Updates(req).Error; err != nil {
		r.log.Error("Ошибка при обновлении отзыва",
			"error", err.Error())
		return fmt.Errorf("ошибка при обновлении отзыва %w", err)
	}

	r.log.Info("Отзыв обновлен")
	return nil
}

func (r *reviewsRepository) Delete(id uint) error {

	if err := r.reviews.Delete(id).Error; err != nil {
		r.log.Error("Ошибка при удалении отзыва",
			"error", err)
		return fmt.Errorf("ошибка при удалении отзыва %w", err)
	}

	r.log.Info("Отзыв удален")

	return nil
}

func (r *reviewsRepository) GetByUserID(userID uint) ([]models.Reviews, error) {
	var reviews []models.Reviews
	if err := r.reviews.Where("user_id = ?", userID).Find(&reviews).Error; err != nil {
		r.log.Error("Ошибка при поиске отзывов",
			"error", err)
		return nil, fmt.Errorf("ошибка при поиске отзывов %w", err)
	}

	r.log.Info("Отзывы получены")
	return reviews, nil
}

func (r *reviewsRepository) GetByCategoryID(categoryID uint) ([]models.Reviews, error) {
	var reviews []models.Reviews

	if err := r.reviews.Where("category_id = ?", categoryID).Find(&reviews).Error; err != nil {
		r.log.Error("Ошибка при поиске отзывов",
			"error", err)
		return nil, fmt.Errorf("ошибка при поиске отзывов %w", err)
	}

	r.log.Info("Отзывы получены")
	return reviews, nil
}
