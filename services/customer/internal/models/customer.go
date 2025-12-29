package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerStatus string

const (
	CustomerStatusActive    CustomerStatus = "active"
	CustomerStatusInactive  CustomerStatus = "inactive"
	CustomerStatusSuspended CustomerStatus = "suspended"
)

type Customer struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null;type:varchar(255)" json:"email"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	Status    CustomerStatus `gorm:"type:enum('active', 'inactive', 'suspended');default:'active'" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// todo: define relationships if any
}

// sets uuid for as an ID type
func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// overrides the table name
func (Customer) TableName() string {
	return "customers"
}
