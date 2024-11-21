package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	FullName    string `json:"full_name" bson:"full_name"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"` // unique
	DateOfBirth string `json:"date_of_birth" bson:"date_of_birth"`
	AvatarURL   string `json:"avatar_url" bson:"avatar_url"`
	Bio         string `json:"bio" bson:"bio"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`               // MongoDB ID
	Username  string             `json:"username" bson:"username" binding:"required"`              // unique
	Email     string             `json:"email" bson:"email" binding:"required,email"` // unique
	Password  string             `json:"password,omitempty" bson:"password,omitempty" binding:"required"` // hashed
	Role      string             `json:"role" bson:"role"`                      // e.g., VIP-0, VIP-1, Admin
	IsActive  bool               `json:"is_active" bson:"is_active"`            // Account status
	Profile   Profile            `json:"profile" bson:"profile"`                // Nested personal info
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`          // Account creation time
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`          // Last update time
}

type UserDTO struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`               // MongoDB ID
	Username  string             `json:"username" bson:"username" binding:"required"`              // unique
	Email     string             `json:"email" bson:"email" binding:"required,email"` // unique
	Role      string             `json:"role" bson:"role"`                      // e.g., VIP-0, VIP-1, Admin
	IsActive  bool               `json:"is_active" bson:"is_active"`            // Account status
	Profile   Profile            `json:"profile" bson:"profile"`                // Nested personal info
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`          // Account creation time
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`          // Last update time
}
