package models

import (
	"time"
)

type UserAuthAction string

const (
	UserAuthActionRegistered       UserAuthAction = "registered"
	UserAuthActionLoggedIn         UserAuthAction = "logged_in"
	UserAuthActionLoggedout        UserAuthAction = "logged_out"
	UserAuthActiionPasswordChanged UserAuthAction = "password_changed"
	UserAuthActionEmailVerfied     UserAuthAction = "email_verified"
	UserAuthActionAccountLocked    UserAuthAction = "account_locked"
	UserAuthActionAccountUnlocked  UserAuthAction = "account_unlocked"
)

type UserAuthAuditLog struct {
	ID          int64          `gorm:"primarykey;AutoIncrement" json:"id"`
	UserID      string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Action      UserAuthAction `gorm:"type:enum('registered', 'logged_in', 'logged_out', 'account_locked', 'account_unlocked');not null;index" json:"action"`
	IPAddress   string         `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent   string         `gorm:"type:text" json:"user_agent"`
	Metadata    JSON           `gorm:"type:json" json:"metadata"`
	PerformedAt time.Time      `gorm:"autoCreateTime;index" json:"performed_at"`
	User        User           `gorm:"foreignkey:UserID;references:ID" json:"user,omitempty"`
}

func (UserAuthAuditLog) TableName() string {
	return "user_auth_audit_logs"
}
