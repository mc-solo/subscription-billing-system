package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email            string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash     string `gorm:"not null" json:"-"`
	FullName         string `json:"full_name"`
	StripeCustomerID string `gorm:"index" json:"stripe_customer_id"`
	Status           string `gorm:"default:'active'" json:"status"` //active, past_due, canceled
}
