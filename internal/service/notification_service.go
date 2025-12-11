package service

import (
	"healthy_body/internal/models"
)

type NotificationService interface {
	SendPaymentSuccess(user *models.User, category *models.Categories) error
}
