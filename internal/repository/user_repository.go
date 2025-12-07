package repository

import (
	"fmt"
	"healthy_body/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(req *models.User) error
	GetAllUser() ([]models.User, error)
	GetUserByID(id uint) (*models.User, error)
	Update(user *models.User, id uint) error
	Delete(id uint) error
}

type gormUserRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewUserRepository(db *gorm.DB, log *slog.Logger) UserRepository {
	return &gormUserRepository{
		db:  db,
		log: log,
	}
}

func (r *gormUserRepository) Create(req *models.User) error {
	if err := r.db.Create(req).Error; err != nil {
		r.log.Error("Ошибка при создании пользователя в слое репозиторий",
			"name", req.Name,
			"balance", req.Balance,
			"error", err.Error(),
		)

		return fmt.Errorf("ошибка при создании пользователя")
	}

	r.log.Info("Пользователь успешно создан",
		"имя", req.Name,
		"балланс", req.Balance,
	)

	return nil
}

func (r *gormUserRepository) GetAllUser() ([]models.User, error) {
	var users []models.User

	if err := r.db.Find(&users).Error; err != nil {
		r.log.Error("Ошибка при выводе пользователей",
			"error", err.Error(),
		)

		return nil, fmt.Errorf("ошибка при выводе пользователей")
	}

	r.log.Info("Пользователи успешно выведены")
	return users, nil

}

func (r *gormUserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, id).Error; err != nil {
		r.log.Error("Ошибка при получении пользователя по ID",
			"id", id,
			"error", err.Error())
		return nil, fmt.Errorf("ошибка при получении пользователя по %d", id)
	}

	r.log.Info("Пользователь найден успешно",
		"id", user.ID,
		"name", user.Name)

	return &user, nil
}

func (r *gormUserRepository) Update(req *models.User, id uint) error {

	newUser := models.User{
		Name:    req.Name,
		Balance: req.Balance,
	}

	result := r.db.Model(&models.User{}).Where("id = ?", id).Updates(&newUser)

	if result.Error != nil {
		r.log.Error("Ошибка при обновлении пользователя по ID",
			"id", id,
			"имя", req.Name,
			"баланс", req.Balance,
			"error", result.Error.Error(),
		)
		return fmt.Errorf("ошибка при обновлении пользователя по %d", id)
	}

	r.log.Info("Пользователь успешно обновлен")
	return nil
}

func (r *gormUserRepository) Delete(id uint) error {

	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		r.log.Error("Ошибка при удалении пользователя",
			"error", err.Error())
		return fmt.Errorf("ошибка при удалении пользователя")
	}

	r.log.Info("Пользователь удален")
	return nil

}
