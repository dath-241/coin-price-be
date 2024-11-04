package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Alert struct {
	ID                 primitive.ObjectID `json:"alert_id" bson:"_id,omitempty"`
	Symbol             string             `json:"symbol"`
	Price              float64            `json:"price"`
	Condition          string             `json:"condition"`
	IsActive           bool               `json:"is_active"`
	NotificationMethod string             `json:"notification_method,omitempty"`
}
