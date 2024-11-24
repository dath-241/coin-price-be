package models

// User đại diện cho một người dùng trong hệ thống
type User struct {
	ID     string  `json:"id"`     // ID của người dùng
	Email  string  `json:"email"`  // Email của người dùng
	Alerts []Alert `json:"alerts"` // Danh sách các cảnh báo của người dùng
}
