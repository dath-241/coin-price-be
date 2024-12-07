package repository

import (
	"context"
	//"fmt"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository interface {
    FindPayments(ctx context.Context, filter bson.M) ([]models.Order, error)
}

type MongoPaymentRepository struct {
    Collection *mongo.Collection
}

func (r *MongoPaymentRepository) FindPayments(ctx context.Context, filter bson.M) ([]models.Order, error) {
    //fmt.Println("Filter used for finding payments:", filter)

    cursor, err := r.Collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var payments []models.Order
    if err := cursor.All(ctx, &payments); err != nil {
        return nil, err
    }

    return payments, nil
}

