package repository

import (
	"context"
	//"fmt"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"

	//"fmt"
	//"test/models"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
)

type MockPaymentRepository struct {
	Payments map[string]interface{} // Dùng map để lưu các đơn thanh toán theo ID
	Err      error
}

func (m *MockPaymentRepository) FindPayments(ctx context.Context, filter bson.M) ([]models.Order, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	var payments []models.Order
	for _, payment := range m.Payments {
		if order, ok := payment.(models.Order); ok {
			payments = append(payments, order)
		}
	}
	return payments, nil
}

// func (m *MockPaymentRepository) FindPayments(ctx context.Context, filter bson.M) ([]models.Order, error) {
// 	if m.Err != nil {
// 		return nil, m.Err
// 	}

// 	var payments []models.Order
// 	// Lọc theo user_id (hoặc các điều kiện khác từ filter)
// 	userID, ok := filter["user_id"].(string)
// 	if !ok || userID == "" {
// 		return nil, fmt.Errorf("invalid user_id filter")
// 	}

// 	for _, payment := range m.Payments {
// 		if order, ok := payment.(models.Order); ok {
// 			if order.UserID.Hex() == userID {
// 				payments = append(payments, order)
// 			}
// 		}
// 	}
// 	return payments, nil
// }

