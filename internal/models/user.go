package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string      `json:"name"`
	Balance      int         `json:"balance"`
	Email        string      `json:"email"`
	CategoriesID uint        `json:"categories_id"`
	Categories   *Categories `json:"-" gorm:"foreignKey:CategoriesID"`

	UserSubscriptions []UserSubscription `json:"userSubscriptions"`
	UserPlans         []UserPlan         `json:"userPlans"`
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
