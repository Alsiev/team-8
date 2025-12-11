package models

import "gorm.io/gorm"

type UserPlan struct {
    gorm.Model
    UserID     uint
    CategoriesID uint

    User     *User     		`gorm:"foreignKey:UserID"`
    Categories *Categories 	`gorm:"foreignKey:CategoriesID"` // обязательно указать foreignKey
}