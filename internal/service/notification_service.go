package service

import (
	"healthy_body/internal/models"
	"log/slog"
)

type NotificationService interface {
	SendPaymentSuccess(user *models.User, category *models.Category) error
}

type notificationService struct {
	logger *slog.Logger
}

func NewNotificationService(logger *slog.Logger) NotificationService {
	return &notificationService{
		logger: logger,
	}
}

func (s *notificationService) SendPaymentSuccess(user *models.User, category *models.Category) error {

	s.logger.Info("Уведомление: успешная оплата",
		"user", user.Name,
		"category", category.Name,
	)

	return nil
}
