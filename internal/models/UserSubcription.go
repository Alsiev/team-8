package models

import (
	"time"

	"gorm.io/gorm"
)

type UserSubscription struct {
	gorm.Model
	UserID         uint      `json:"user_id"`
	SubscriptionID uint      `json:"subscription_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	IsActive       bool      `json:"is_active"`

	User         *User         `json:"-"`
	Subscription *Subscription `json:"-"`
}