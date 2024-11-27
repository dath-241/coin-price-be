package models

import (
    // "time"
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
    CreatedAt  primitive.DateTime          `bson:"created_at"`  // Ngày tạo đơn hàng
    UpdatedAt  primitive.DateTime          `bson:"updated_at"`  // Ngày cập nhật đơn hàng
    TransactionStatus string      `bson:"transaction_status"` // Thêm trường này để lưu trạng thái giao dịch
}

type CreateVIPPaymentRequest struct {
    Amount   int    `json:"amount"`
    VIPLevel string `json:"vip_level"`
}
type CreateVIPPaymentReponse struct {
	PaymentURL string `json:"payment_url"`
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