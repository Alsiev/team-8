package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name       string    `json:"name"`
	Balance    int       `json:"balance"`
	Email      string    `json:"email"`
	CategoryID uint      `json:"category_id"`
	Category   *Category `json:"-" gorm:"foreignKey:CategoryID"`

	UserSubscriptions []UserSubscription `json:"-"`
	UserPlans         []UserPlan         `json:"-"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name    *string `json:"name"`
	Balance *int    `json:"balance"`
	Email   *string `json:"email"`
}
