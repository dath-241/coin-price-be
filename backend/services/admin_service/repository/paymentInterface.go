package repository

import (
	"context"
	"fmt"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository interface {
    FindPayments(ctx context.Context, filter bson.M) ([]models.Order, error)
}

type MongoPaymentRepository struct {
    Collection *mongo.Collection
}


func (r *MongoPaymentRepository) FindPayments(ctx context.Context, filter bson.M) ([]models.Order, error) {
    fmt.Println("Filter used for finding payments:", filter) // Log filter để xem nó có đúng không

    cursor, err := r.Collection.Find(ctx, filter)
    if err != nil {
        // Log lỗi chi tiết khi gặp lỗi trong quá trình Find
        fmt.Println("Error during Find operation:", err)
        return nil, err
    }
    defer cursor.Close(ctx)

    var payments []models.Order
    if err := cursor.All(ctx, &payments); err != nil {
        // Log lỗi khi giải mã dữ liệu trả về
        fmt.Println("Error during cursor.All operation:", err)
        return nil, err
    }

    // Log số lượng payments tìm được để xác nhận kết quả
    fmt.Println("Number of payments found:", len(payments))
    return payments, nil
}

