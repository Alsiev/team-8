package service

import (
	"fmt"
	"healthy_body/internal/models"
	"log/slog"

	gomail "gopkg.in/gomail.v2"
)

type EmailNotificationService struct {
	fromEmail string
	fromPass  string
	smtpHost  string
	smtpPort  int
	logger    *slog.Logger
}

func NewEmailNotificationService(
	fromEmail, fromPass, smtpHost string,
	smtpPort int,
	logger *slog.Logger,
) *EmailNotificationService {

	return &EmailNotificationService{
		fromEmail: fromEmail,
		fromPass:  fromPass,
		smtpHost:  smtpHost,
		smtpPort:  smtpPort,
		logger:    logger,
	}
}

func (s *EmailNotificationService) SendPaymentSuccess(user *models.User, category *models.Category) error {

	if user.Email == "" {
		s.logger.Warn("у пользователя нет email", "user_id", user.ID)
		return nil
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", s.fromEmail)
	msg.SetHeader("To", user.Email)
	msg.SetHeader("Subject", "Оплата прошла успешно")
	msg.SetBody("text/plain",
		fmt.Sprintf(
			"Привет, %s!\n\nВы успешно оплатили категорию: %s.\nСпасибо, что пользуетесь нашим сервисом!",
			user.Name,
			category.Name,
		),
	)

	dialer := gomail.NewDialer(s.smtpHost, s.smtpPort, s.fromEmail, s.fromPass)

	dialer.SSL = false

	if err := dialer.DialAndSend(msg); err != nil {
		s.logger.Error("не удалось отправить email", "err", err)
		return err
	}

	s.logger.Info("email уведомление отправлено",
		"to", user.Email,
		"category", category.Name,
	)

	return nil
}
