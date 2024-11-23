package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // Trường ID tự động tạo bởi MongoDB
	Name      string             `json:"name"`
	Email     string             `json:"email" binding:"required,email"`
	Password  string             `json:"password,omitempty" binding:"required"` // Không trả về password trong response
	Role      string             `json:"role"`
	CreatedAt time.Time          `json:"created_at"`
}
