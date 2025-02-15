package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	FullName    string `json:"full_name" bson:"full_name" example:"John Doe"`
	PhoneNumber string `json:"phone_number" bson:"phone_number" example:"+84911123456"` // unique
	DateOfBirth string `json:"date_of_birth" bson:"date_of_birth" example:"1995-05-15"`
	AvatarURL   string `json:"avatar_url" bson:"avatar_url" example:"https://example.com/avatar.png"`
	Bio         string `json:"bio" bson:"bio" example:"Software developer based in NY"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"648fd7ef2cbae153e4b5c7df"`                                     // MongoDB ID
	Username  string             `json:"username" bson:"username" binding:"required" example:"johndoe"`                                  // unique
	Email     string             `json:"email" bson:"email" binding:"required,email" example:"user@example.com"`                         // unique
	Password  string             `json:"password,omitempty" bson:"password,omitempty" binding:"required" example:"hashed_password_here"` // hashed
	Role      string             `json:"role" bson:"role" example:"VIP-0"`                                                               // e.g., VIP-0, VIP-1, Admin
	IsActive  bool               `json:"is_active" bson:"is_active"`                                                                     // Account status
	Profile   Profile            `json:"profile" bson:"profile"`                                                                         // Nested personal info
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at" example:"2024-11-01T10:00:00Z"`                                    // Account creation time
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at" example:"2024-11-01T10:00:00Z"`
	Alerts    []string           `json:"alerts"` // Last update time
}

type UserDTO struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"648fd7ef2cbae153e4b5c7df"`             // MongoDB ID
	Username  string             `json:"username" bson:"username" binding:"required" example:"johndoe"`          // unique
	Email     string             `json:"email" bson:"email" binding:"required,email" example:"user@example.com"` // unique
	Role      string             `json:"role" bson:"role" example:"VIP-0"`                                       // e.g., VIP-0, VIP-1, Admin
	IsActive  bool               `json:"is_active" bson:"is_active"`                                             // Account status
	Profile   Profile            `json:"profile" bson:"profile"`                                                 // Nested personal info
	CreatedAt time.Time          `json:"created_at" bson:"created_at" example:"2024-11-01T10:00:00Z"`            // Account creation time
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at" example:"2024-11-01T10:00:00Z"`            // Last update time
}

type ErrorResponse struct {
	Error string `json:"error" example:"An error occurred"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UserResponse  struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    VIPLevel string `json:"vip_level"`
    Status   bool   `json:"status"`
}


type RegisterRequest struct {
	Username  string      `json:"username" bson:"username" binding:"required" example:"johndoe"`              // unique
	Email     string      `json:"email" bson:"email" binding:"required,email" example:"user@example.com"` // unique
	Password  string      `json:"password,omitempty" bson:"password,omitempty" binding:"required" example:"hashed_password_here"`
	Profile   Profile     `json:"profile" bson:"profile"`                // Nested personal info
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RorLResponse struct {
	Message string `json:"message"`
	Token 	string `json:"token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetRequest struct {
	OTP       	string `json:"otp" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdateUserProfileRequest struct {
    Name        string `json:"name" example:"John Doe"`
    Username    string `json:"username" example:"johndoe123"`
    Phone       string `json:"phone" example:"+1234567890"`
    Avatar      string `json:"avatar" example:"https://example.com/avatar.jpg"`
    Bio         string `json:"bio" example:"Software Engineer"`
    DateOfBirth string `json:"dateOfBirth" example:"2000-01-01"` // Định dạng: YYYY-MM-DD
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type ChangeMailRequest struct {
	Email string `json:"email" binding:"required,email"` // email mới cần cập nhật
}
