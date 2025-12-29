package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubcriptionStatus string

const (
	SubcriptionStatusActive    SubcriptionStatus = "active"
	SubscriptionStatusCanceled SubcriptionStatus = "canceled"
	SubscriptionStatusPastDue  SubcriptionStatus = "past_due"
	SubscriptionStatusUnpaid   SubcriptionStatus = "unpaid"
)

type Subscription struct {
	ID                   string            `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CustomerId           string            `gorm:"not null;type:varchar(36);index" json:"customer_id"`
	PlanId               string            `gorm:"not null;type:varchar(50)" json:"plan_id"`
	Status               SubcriptionStatus `gorm:"type:enum('active', 'canceled', 'past_due', 'unpaid');default('active')" json:"status"`
	CurrentPeriodStart   time.Time         `gorm:"not null" json:"current_period_start"`
	CurrentPeriodEnd     time.Time         `gorm:"not null" json:"current_period_end"`
	CancelAtPeriodEnd    bool              `gorm:"default:false" json:"cancel_at_period_end"`
	StripeSubscriptionId string            `gorm:"uniqueIndex;type:varchar(255);null" json:"stripe_subscription_id"`
	Metadata             JSON              `gorm:"type:json" json:"metadata"`
	CreatedAt            time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt    `gorm:"index" json:"deleted_at,omitempty"`

	// relations
	Customer Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
}

type JSON map[string]interface{}

func (j JSON) Value() (interface{}, error) {
	return json.Marshal(j)
}

// TODO: unmarshal the raw data to our map
// BeforeCreate hook

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}

	// set default metadata if empty
	if s.Metadata == nil {
		s.Metadata = JSON{}
	}

	return nil
}

// override tablename
func (Subscription) TableName() string {
	return "subscriptions"
}
