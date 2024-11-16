package models

import (
	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims là cấu trúc lưu các claims có thể có trong JWT, ví dụ userID và role.
type CustomClaims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}