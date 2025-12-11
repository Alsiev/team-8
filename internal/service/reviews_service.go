package service

import (
	"fmt"
	"healthy_body/internal/models"
	"healthy_body/internal/repository"
	"log/slog"
)

type ReviewsService interface {
	CreateReview(req models.CreateReviewRequest, userID uint) (uint, error)
	GetReview(id uint) (*models.GetReview, error)
	GetReviewsByUser(userID uint) ([]models.GetReview, error)
	GetReviewsByCategory(categoryID uint) ([]models.GetReview, error)
	UpdateReview(id uint, req models.UpdateReviewRequest, userID uint) error
	DeleteReview(id uint, userID uint) error
}

type reviewsService struct {
	repo repository.ReviewsRepository
	log  *slog.Logger
}

func NewReviewsService(repo repository.ReviewsRepository, log *slog.Logger) ReviewsService {
	return &reviewsService{repo: repo, log: log}
}

func (s *reviewsService) CreateReview(req models.CreateReviewRequest, userID uint) (uint, error) {

	if userID == 0 {
		s.log.Warn("Такого пользователя не существует",
			"user_id", userID)
		return 0, fmt.Errorf("такого пользователя не существует")
	}

	if req.CategoriesID == 0 {
		s.log.Warn("Такой категории нету",
			"category_id", req.CategoriesID)
		return 0, fmt.Errorf("такой категории не существует")

	}

	if req.Rating < 0 || req.Rating > 5 {
		s.log.Warn("Оценка должна быть выбрана от 1 до 5",
			"ваша оценка", req.Rating)
		return 0, fmt.Errorf("оценка должна быть выбрана от 1 до 5")
	}

	newReview := models.Reviews{
		UserID:     req.UserID,
		CategoriesID: req.CategoriesID,
		Rating:     req.Rating,
		Content:    req.Content,
	}

	if err := s.repo.CreateReviews(&newReview); err != nil {
		s.log.Error("Ошибка при создании отзыва",
			"error", err.Error())
		return 0, fmt.Errorf("ошибка при создании отзыва")
	}

	return newReview.ID, nil

}

func (s *reviewsService) GetReview(id uint) (*models.GetReview, error) {
	if id == 0 {
		s.log.Warn("id не указан")
		return nil, fmt.Errorf("id не указан")
	}

	req, err := s.repo.GetReviewsByID(id)
	if err != nil {
		s.log.Error("Ошибка при выводе отзыва",
			"id", id,
			"error", err.Error())
		return nil, fmt.Errorf("ошибка при выводе отзыва: %w", err)
	}

	getReview := &models.GetReview{
		ID:         req.ID,
		CategoriesID: req.CategoriesID,
		UserID:     req.UserID,
		Rating:     req.Rating,
		Content:    req.Content,
		Date:       req.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	s.log.Info("Отзыв получен")

	return getReview, nil
}

func (s *reviewsService) GetReviewsByUser(userID uint) ([]models.GetReview, error) {
	if userID == 0 {
		s.log.Warn("ID пользователя не указан")
		return nil, fmt.Errorf("ID пользователя не указан")
	}

	reviews, err := s.repo.GetByUserID(userID)
	if err != nil {
		s.log.Error("Ошибка при получении отзывов пользователя",
			"user_id", userID,
			"error", err.Error())
		return nil, fmt.Errorf("ошибка при получении отзывов пользователя: %w", err)
	}

	var result []models.GetReview
	for _, review := range reviews {
		getReview := models.GetReview{
			ID:         review.ID,
			CategoriesID: review.CategoriesID,
			UserID:     review.UserID,
			Rating:     review.Rating,
			Content:    review.Content,
			Date:       review.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		result = append(result, getReview)
	}

	s.log.Info("Отзывы пользователя получены",
		"user_id", userID,
		"count", len(result))
	return result, nil
}

func (s *reviewsService) GetReviewsByCategory(categoryID uint) ([]models.GetReview, error) {
	if categoryID == 0 {
		s.log.Warn("ID категории не указан")
		return nil, fmt.Errorf("ID категории не указан")
	}

	reviews, err := s.repo.GetByCategoryID(categoryID)
	if err != nil {
		s.log.Error("Ошибка при получении отзывов по категории",
			"category_id", categoryID,
			"error", err.Error())
		return nil, fmt.Errorf("ошибка при получении отзывов по категории: %w", err)
	}

	var result []models.GetReview
	for _, review := range reviews {
		getReview := models.GetReview{
			ID:         review.ID,
			CategoriesID: review.CategoriesID,
			UserID:     review.UserID,
			Rating:     review.Rating,
			Content:    review.Content,
			Date:       review.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		result = append(result, getReview)
	}

	s.log.Info("Отзывы по категории получены",
		"category_id", categoryID,
		"count", len(result))
	return result, nil
}

func (s *reviewsService) UpdateReview(id uint, req models.UpdateReviewRequest, userID uint) error {
	if id == 0 {
		s.log.Warn("ID отзыва не указан")
		return fmt.Errorf("ID отзыва не указан")
	}

	if userID == 0 {
		s.log.Warn("ID пользователя не указан")
		return fmt.Errorf("ID пользователя не указан")
	}

	review, err := s.repo.GetReviewsByID(id)
	if err != nil {
		s.log.Error("Отзыв не найден",
			"id", id,
			"error", err.Error())
		return fmt.Errorf("отзыв не найден: %w", err)
	}

	if review.UserID != userID {
		s.log.Warn("Попытка обновления чужого отзыва",
			"user_id", userID,
			"review_user_id", review.UserID)
		return fmt.Errorf("нельзя обновлять чужой отзыв")
	}

	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 5 {
			s.log.Warn("Оценка должна быть от 1 до 5",
				"rating", *req.Rating)
			return fmt.Errorf("оценка должна быть от 1 до 5")
		}
		review.Rating = *req.Rating
	}

	if req.Content != nil {
		review.Content = *req.Content
	}

	if err := s.repo.UpdateReviews(review); err != nil {
		s.log.Error("Ошибка при обновлении отзыва",
			"id", id,
			"error", err.Error())
		return fmt.Errorf("ошибка при обновлении отзыва: %w", err)
	}

	s.log.Info("Отзыв успешно обновлен",
		"id", id)
	return nil
}

func (s *reviewsService) DeleteReview(id uint, userID uint) error {
	if id == 0 {
		s.log.Warn("ID отзыва не указан")
		return fmt.Errorf("ID отзыва не указан")
	}

	if userID == 0 {
		s.log.Warn("ID пользователя не указан")
		return fmt.Errorf("ID пользователя не указан")
	}

	review, err := s.repo.GetReviewsByID(id)
	if err != nil {
		s.log.Error("Отзыв не найден",
			"id", id,
			"error", err.Error())
		return fmt.Errorf("отзыв не найден: %w", err)
	}

	if review.UserID != userID {
		s.log.Warn("Попытка удаления чужого отзыва",
			"user_id", userID,
			"review_user_id", review.UserID)
		return fmt.Errorf("нельзя удалять чужой отзыв")
	}

	if err := s.repo.Delete(id); err != nil {
		s.log.Error("Ошибка при удалении отзыва",
			"id", id,
			"error", err.Error())
		return fmt.Errorf("ошибка при удалении отзыва: %w", err)
	}

	s.log.Info("Отзыв успешно удален",
		"id", id)
	return nil
}
