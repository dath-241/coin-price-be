package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Alert represents a price alert or symbol alert in the system
type Alert struct {
	ID                 primitive.ObjectID `json:"alert_id" bson:"_id,omitempty"`
	UserID             string             `json:"user_id" bson:"user_id"`                                         // ID của người dùng đặt cảnh báo
	Symbol             string             `json:"symbol" bson:"symbol"`                                           // Symbol for the alert, e.g., "BTCUSDT"
	Price              float64            `json:"price,omitempty" bson:"price,omitempty"`                         // Price threshold for the alert, only applicable for price alerts
	Condition          string             `json:"condition,omitempty" bson:"condition,omitempty"`                 // Condition for the alert, e.g., ">=", "<=", "==", "in range", "out range"
	Threshold          float64            `json:"threshold,omitempty" bson:"threshold,omitempty"`     // Threshold for the alert, only applicable for price alerts
	Market             string             `json:"market,omitempty" bson:"market,omitempty"`                       // Market for the alert, e.g., "spot" or "future"
	Range              []float64          `json:"range,omitempty" bson:"range,omitempty"`                         // Price range for alerts that check if price falls within or outside a range
	IsActive           bool               `json:"is_active" bson:"is_active"`                                     // Whether the alert is active
	NotificationMethod string             `json:"notification_method" bson:"notification_method"`                 // How the user will be notified (email, push, Telegram)
	Type               string             `json:"type,omitempty" bson:"type,omitempty"`                           // Type of symbol alert, e.g., "new_listing" or "delisting"
	Symbols            []string           `json:"symbols,omitempty" bson:"symbols,omitempty"`                     // List of symbols for symbol alerts
	Frequency          string             `json:"frequency,omitempty" bson:"frequency,omitempty"`                 // Frequency of notification, e.g., "immediate", "daily", "weekly"
	CreatedAt          primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`               // Timestamp for when the alert was created
	UpdatedAt          primitive.DateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`               // Timestamp for when the alert was last updated
	SnoozeCondition    string             `json:"snooze_condition,omitempty" bson:"snooze_condition,omitempty"`   // Loại snooze
	MaxRepeatCount     int                `json:"max_repeat_count,omitempty" bson:"max_repeat_count,omitempty"`   // Số lần lặp lại tối đa
	NextTriggerTime    time.Time          `json:"next_trigger_time,omitempty" bson:"next_trigger_time,omitempty"` // Thời gian kích hoạt tiếp theo
	RepeatCount        int                `json:"repeat_count,omitempty" bson:"repeat_count,omitempty"`           // Số lần đã lặp lại
	Message            string             `json:"message" bson:"message"`                                         // Thông điệp sẽ được gửi trong cảnh báo (ví dụ: "BTC giá đã vượt $20,000")
	LastInterval       string             `json:"last_fundingrate_interval" bson:"last_fundingrate_interval"`
}


// Symbol struct
type Symbol struct {
	Symbol string `json:"symbol"`
	Status string `json:"status"`
}

// PriceDifferenceAlert struct
type PriceDifferenceAlert struct {
	ID          string  `json:"id"`
	SpotPrice   float64 `json:"spot_price"`
	FuturePrice float64 `json:"future_price"`
	Threshold   float64 `json:"threshold"`
	Triggered   bool    `json:"triggered"`
	Action      string  `json:"action"`
}

// NewAlert creates a new alert with default values, including checking Condition for range requirements
func NewAlert(userID, symbol, market, condition, notificationMethod string, price float64, alertType, frequency, snoozeCondition string, maxRepeatCount int, symbols []string, priceRange []float64) *Alert {
	var finalRange []float64
	if condition == "in range" || condition == "out range" {
		finalRange = priceRange // Set the range only for in/out range conditions
	} else {
		finalRange = nil // Default to nil if range is not applicable
	}

	return &Alert{
		ID:                 primitive.NewObjectID(),
		UserID:             userID,
		Symbol:             symbol,
		Market:             market,
		Price:              price,
		Condition:          condition,
		IsActive:           true,
		NotificationMethod: notificationMethod,
		Type:               alertType,
		Symbols:            symbols,
		Frequency:          frequency,
		SnoozeCondition:    snoozeCondition,
		MaxRepeatCount:     maxRepeatCount,
		Range:              finalRange,
		CreatedAt:          primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:          primitive.NewDateTimeFromTime(time.Now()),
		NextTriggerTime:    time.Now().Add(time.Minute),
		RepeatCount:        0,
	}
}
