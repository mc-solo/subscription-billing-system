package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusPendingVerification UserStatus = "pending_verification"
	UserStatusActice              UserStatus = "active"
	UserStatusSuspended           UserStatus = "suspended"
)

type User struct {
	ID                          string     `gorm:";type:varchar(36)" json:"id"`
	Email                       string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash                string     `gorm:"not null" json:"-"`
	FullName                    string     `json:"full_name"`
	StripeCustomerID            string     `gorm:"index" json:"stripe_customer_id"`
	Status                      UserStatus `gorm:"default:'active'" json:"status"` //active, past_due, canceled
	EmailVerified               bool       `gorm:"default:false;type:varchar(255)" json:"email_verfied"`
	VerificationToken           *string    `gorm:"type:varchar(255);null" json:"-"`
	VerificationTokenExipresAt  *time.Time `gorm:"null" json:"-"`
	PasswordResetToken          *string    `gorm:"type:varchar(100);null" json:"-"`
	PasswordResetTokenExpiresAt *time.Time `gorm:"null" json:"-"`
	LastLoginAt                 *time.Time `gorm:"null" json:"last_login_at,omitempty"`
	FailedLoginAttempts         int        `gorm:"default:0" json:"-"`
	LockedUntill                *time.Time `gorm:"null" json:"-"`

	// time stamps
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// relationships
	Customers []Customer `gorm:"foreignkey:UserID;references:ID" json:"customers,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (User) TableName() string {
	return "users"
}
