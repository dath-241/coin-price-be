package models

import (
	models "github.com/dath-241/coin-price-be-go/services/trigger-service/models/alert"
)

// User đại diện cho một người dùng trong hệ thống
type User struct {
	ID     string         `json:"id"`     // ID của người dùng
	Email  string         `json:"email"`  // Email của người dùng
	Alerts []models.Alert `json:"alerts"` // Danh sách các cảnh báo của người dùng
}
