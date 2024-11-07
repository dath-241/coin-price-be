package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Alert represents a price alert or symbol alert in the system
type Alert struct {
	ID                 primitive.ObjectID `json:"alert_id" bson:"_id,omitempty"`
	Symbol             string             `json:"symbol" bson:"symbol"`                             // Symbol for the alert, e.g., "BTCUSDT"
	Price              float64            `json:"price,omitempty" bson:"price,omitempty"`           // Price threshold for the alert, only applicable for price alerts
	Condition          string             `json:"condition,omitempty" bson:"condition,omitempty"`   // Condition for the alert, e.g., ">=", "<=", "=="
	IsActive           bool               `json:"is_active" bson:"is_active"`                       // Whether the alert is active
	NotificationMethod string             `json:"notification_method" bson:"notification_method"`   // How the user will be notified (email, push, Telegram)
	Type               string             `json:"type,omitempty" bson:"type,omitempty"`             // Type of symbol alert, e.g., "new_listing" or "delisting"
	Symbols            []string           `json:"symbols,omitempty" bson:"symbols,omitempty"`       // List of symbols for symbol alerts
	Frequency          string             `json:"frequency,omitempty" bson:"frequency,omitempty"`   // Frequency of notification, e.g., "immediate", "daily", "weekly"
	CreatedAt          primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"` // Timestamp for when the alert was created
	UpdatedAt          primitive.DateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"` // Timestamp for when the alert was last updated
}

type Symbol struct {
    Symbol string `json:"symbol"`
    Status string `json:"status"`
}

// NewAlert creates a new alert with default values
func NewAlert(symbol, condition, notificationMethod string, price float64, alertType, frequency string, symbols []string) *Alert {
	return &Alert{
		ID:                 primitive.NewObjectID(),
		Symbol:             symbol,
		Price:              price,
		Condition:          condition,
		IsActive:           true,
		NotificationMethod: notificationMethod,
		Type:               alertType,
		Symbols:            symbols,
		Frequency:          frequency,
		CreatedAt:          primitive.NewDateTimeFromTime(time.Now()),
	}
}
