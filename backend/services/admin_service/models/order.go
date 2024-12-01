package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
    ID         primitive.ObjectID `bson:"_id,omitempty"`
    UserID     primitive.ObjectID `bson:"user_id"`     // ID của người dùng
    VipLevel   string             `bson:"vip_level"`   // Cấp VIP
    Amount     int                `bson:"amount"`      // Số tiền thanh toán
    OrderID    string             `bson:"order_id"`    // Mã đơn hàng
    OrderInfo  string             `bson:"orderInfo"`    // Thông tin đơn hàng
    PaymentURL string             `bson:"payment_url"` // URL thanh toán MoMo
    CreatedAt  primitive.DateTime `bson:"created_at"`  // Ngày tạo đơn hàng
    UpdatedAt  primitive.DateTime `bson:"updated_at"`  // Ngày cập nhật đơn hàng
    TransactionStatus string      `bson:"transaction_status"` // Thêm trường này để lưu trạng thái giao dịch
}

type CreateVIPPaymentRequest struct {
    Amount   int    `json:"amount"`
    VIPLevel string `json:"vip_level"`
}
type CreateVIPPaymentReponse struct {
	PaymentURL  string `json:"payment_url"`
	OrderID 	string `json:"order_id"`
}

type QueryPaymentRequest struct {
    OrderID   string `json:"orderId"`
    RequestID string `json:"requestId"`
    Lang      string `json:"lang"`
}
type ReponseQueryPaymentRequest struct {
    Message   string `json:"message"`
    Status    string `json:"status"`
}

type MoMoResponse struct {
	PartnerCode  string `json:"partnerCode"`
	RequestId    string `json:"requestId"`
	OrderId      string `json:"orderId"`
	ResultCode   string `json:"resultCode"`
	Message      string `json:"message"`
	ResponseTime string `json:"responseTime"`
	ExtraData    string `json:"extraData"`
	Signature    string `json:"signature"`
}

// PaymentDetails defines the structure for payment history details
type PaymentDetailsUser struct {
    OrderInfo         string    `json:"order_info"`         // Thông tin đơn hàng
    TransactionStatus string    `json:"transaction_status"` // Trạng thái giao dịch
    Amount            float64   `json:"amount"`             // Số tiền thanh toán
    CreatedAt         time.Time `json:"created_at"`         // Thời gian tạo
    UpdatedAt         time.Time `json:"updated_at"`         // Thời gian cập nhật
}

// PaymentHistoryResponse is the response structure for payment history
type PaymentHistoryUserResponse struct {
    PaymentHistory []PaymentDetailsUser `json:"payment_history"` // Danh sách lịch sử thanh toán
}

type PaymentAdmin struct {
    OrderID             string    `json:"order_id"`         // Thông tin đơn hàng
    UserID              string    `bson:"user_id"` 
    OrderInfo           string    `json:"order_info"` 
    TransactionStatus   string    `json:"transaction_status"` // Trạng thái giao dịch
    Amount              float64   `json:"amount"`             // Số tiền thanh toán
}

// PaymentHistoryResponse is the response structure for payment history
type PaymentHistoryAdminResponse struct {
    PaymentHistory []PaymentAdmin `json:"payment_history"` // Danh sách lịch sử thanh toán
}

type PaymentDetailsAdmin struct {
    OrderID             string      `json:"order_id"`         // Thông tin đơn hàng
    OrderInfo           string      `json:"order_info"` 
    TransactionStatus   string      `json:"transaction_status"` // Trạng thái giao dịch
    Amount              float64     `json:"amount"`             // Số tiền thanh toán
    VipLevel            string      `json:"vip_level"`
    PaymentURL          string      `json:"payment_url"` // URL thanh toán MoMo
    CreatedAt           time.Time   `json:"created_at"`  // Ngày tạo đơn hàng
    UpdatedAt           time.Time   `json:"updated_at"` 
}

type PaymentHisDetailsAdminResponse struct {
    PaymentHistory []PaymentDetailsAdmin `json:"payment_history"` // Danh sách lịch sử thanh toán
}