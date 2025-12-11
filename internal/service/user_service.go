package service

import (
	"fmt"
	"healthy_body/internal/models"
	"healthy_body/internal/repository"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(req models.CreateUserRequest) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserPlan(userID uint) (*models.Category, error)
	UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error)
	Delete(id uint) error

	Payment(userID uint, categoryID uint) error
	SubPayment(userID, subID uint) error
}

type userService struct {
	userRepo repository.UserRepository
	log      *slog.Logger
	db       *gorm.DB
	sub      SubscriptionService
	categoryRepo repository.CategoryRepo
}

func NewUserService(userRepo repository.UserRepository, log *slog.Logger, db *gorm.DB, sub SubscriptionService, categoryRepo repository.CategoryRepo) UserService {
	return &userService{
		userRepo: userRepo,
		log:      log,
		db:       db,
		sub:      sub,
		categoryRepo: categoryRepo,
	}
}

func (s *userService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	if len(req.Name) < 2 {
		s.log.Warn("имя пользователя слишком короткое",
			"имя", req.Name,
			"минимальное количество символов имени", 2,
		)
		return nil, fmt.Errorf("имя должно содержать минимум 2 символа")

	}

	newUser := &models.User{
		Name:       req.Name,
		Balance:    0,
		CategoryID: 2,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		s.log.Error("Ошибка при создании пользователя",
			"имя", req.Name,
			"error", err.Error())

		return nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	s.log.Info("Пользователь создан",
		"id", newUser.ID,
		"имя", req.Name,
		"баланс", 0,
		"категория", 0)

	return newUser, nil

}

func (s *userService) GetAllUsers() ([]models.User, error) {

	result, err := s.userRepo.GetAllUser()
	if err != nil {
		s.log.Error("Ошибка при выводе пользователей",
			"error", err.Error())
		return nil, fmt.Errorf("ошибка при выводе пользователей: %w", err)
	}

	s.log.Info("Пользователи получены",
		"количество пользователей", len(result))

	return result, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {

	if id == 0 {
		s.log.Warn("id не указан")
		return nil, fmt.Errorf("id не указан")
	}

	result, err := s.userRepo.GetUserByID(id)

	if err != nil {
		s.log.Error("Ошибка при выводе пользователя",
			"id", id,
			"error", err.Error())
		return nil, fmt.Errorf("ошибка при выводе пользователя: %w", err)
	}

	s.log.Info("Пользователь найден",
		"id", result.ID,
		"имя", result.Name,
		"баланс", result.Balance,
	)

	return result, nil
}

func (s *userService) GetUserPlan(userID uint) (*models.Category, error){
 
    user, err := s.userRepo.GetUserByID(userID)
    if err != nil {
        return nil, err
    }
    category, err := s.categoryRepo.GetWithPlans(user.CategoryID)
    if err != nil {
        return nil, err
    }

    return category, nil
}


func (s *userService) UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error) {

	if req.Name == nil && req.Balance == nil {
		s.log.Warn("Нет полей для обновления", "id", id)
		return nil, fmt.Errorf("не указаны поля для обновления")
	}

	if req.Name != nil {
		if len(*req.Name) < 2 {
			s.log.Warn("Короткое имя при обновлении",
				"id", id,
				"name", *req.Name)
			return nil, fmt.Errorf("имя должно содержать минимум 2 символа")
		}
	}

	if req.Balance != nil {
		balance := *req.Balance
		if balance < 0 {
			s.log.Warn("Отрицательный баланс при обновлении",
				"id", id,
				"balance", balance,
			)
			return nil, fmt.Errorf("баланс не может быть отрицательным")
		}
	}

	user, err := s.GetUserByID(id)

	if err != nil {
		s.log.Error("Ошибка при поиске пользователя",
			"error", err.Error())
		return nil, fmt.Errorf("ошибка при поиске пользователя %w", err)
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Balance != nil {
		user.Balance = *req.Balance
	}

	if err := s.userRepo.Update(user); err != nil {
		s.log.Error("Ошибка при обновлении пользователя",
			"error", err.Error())
		return nil, fmt.Errorf("ошибка при обновлении пользователя %w", err)
	}

	s.log.Info("Пользователь обновлен",
		"id", id,
		"имя", req.Name,
		"баланс", req.Balance,
	)
	return user, nil
}

func (s *userService) Delete(id uint) error {

	if err := s.userRepo.Delete(id); err != nil {
		s.log.Error("Ошибка при удалении пользователя",
			"ID", id,
			"error", err)
		return fmt.Errorf("ошибка при удалении пользователя %w", err)
	}

	return nil
}

func (s *userService) Payment(userID uint, categoryID uint) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {

		var user models.User

		if err := tx.First(&user, userID).Error; err != nil {
			s.log.Error("Ошибка при поиске пользователя",
				"error", err.Error())
			return fmt.Errorf("ошибка при поиске пользователя %w", err)
		}

		var category models.Category

		if err := tx.First(&category, categoryID).Error; err != nil {
			s.log.Error("Ошибка при поиске категории",
				"error", err.Error())
			return fmt.Errorf("ошибка при поиске категории %w", err)
		}

		if user.Balance < category.Price {
			s.log.Warn("Недостаточно средств на счету")
			return fmt.Errorf("недостаточно средств на счету")
		}

		user.Balance -= category.Price
		user.CategoryID = categoryID

		userPlan := &models.UserPlan{
			UserID: userID,
			CategoryID: categoryID,
		}

		if err := tx.Create(&userPlan).Error; err != nil {
			s.log.Error("Ошибка при записи покупки пользователя",
				"error", err.Error())
			return fmt.Errorf("ошибка при записи покупки пользователя %w", err)
		}

		if err := tx.Save(&user).Error; err != nil {
			s.log.Error("Ошибка при сохранении пользователя",
				"error", err.Error())
			return fmt.Errorf("ошибка при сохранении пользователя %w", err)
		}

		s.log.Info("Оплата прошла успешно")

		return nil
	})
	return err
}

func (s *userService) SubPayment(userID, subID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {

		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		sub, err := s.sub.GetSubByID(subID)
		if err != nil {
			return fmt.Errorf("subscription not found: %w", err)
		}

		if user.Balance < sub.Price {
			return fmt.Errorf("недостаточно средств")
		}

		user.Balance -= sub.Price

		if err := tx.Save(&user).Error; err != nil {
			return fmt.Errorf("cannot update user balance: %w", err)
		}

		userSub := &models.UserSubscription{
			UserID:         userID,
			SubscriptionID: subID,
			StartDate:      time.Now(),
			EndDate:        time.Now().Add(time.Hour * 24 * time.Duration(sub.DurationDays)),
			IsActive:       true,
		}

		if err := tx.Create(&userSub).Error; err != nil {
			return fmt.Errorf("cannot create user subscription: %w", err)
		}

		return nil
	})
}
