package models

import "gorm.io/gorm"

type UserPlan struct {
	gorm.Model
	UserID uint `json:"user_id"`
	CategoryID uint `json:"category_id"`

	User     *User     `json:"-" gorm:"foreignKey:UserID"`
	Category *Category `json:"-" gorm:"foreignKey:CategoryID"`
}