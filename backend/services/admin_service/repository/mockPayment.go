package repository

import (
	"context"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"go.mongodb.org/mongo-driver/bson"
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
