package models

import "gorm.io/gorm"

type Reviews struct {
	gorm.Model

	CategoriesID uint        `json:"categories_id"`
	Categories   *Categories `json:"-"`
	UserID       uint        `json:"-"`
	User         *User       `json:"-" gorm:"foreignKey:UserID"`
	Rating       int         `json:"-"`
	Content      string      `json:"-" `
}

type GetReview struct {
	ID           uint   `json:"id"`
	CategoriesID uint   `json:"categories_id"`
	UserID       uint   `json:"user_id"`
	Rating       int    `json:"rating"`
	Content      string `json:"content"`
	Date         string `json:"date"`
}

type CreateReviewRequest struct {
	CategoriesID uint   `json:"categories_id"`
	UserID       uint   `json:"user_id"`
	Rating       int    `json:"rating"`
	Content      string `json:"content"`
}

type UpdateReviewRequest struct {
	Rating  *int    `json:"rating,omitempty"`
	Content *string `json:"content,omitempty"`
}
