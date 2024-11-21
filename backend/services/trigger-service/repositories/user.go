package repositories

import (
	"context"
	"fmt"
	"time"

	config "github.com/dath-241/coin-price-be-go/services/admin_service/config"
	alert "github.com/dath-241/coin-price-be-go/services/trigger-service/models/alert"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByID(userID string) (models.User, error) {
	// Lấy collection từ database
	collection := config.DB.Collection("User")

	// Thiết lập timeout cho context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tạo biến lưu trữ kết quả người dùng
	var user models.User

	// Chuyển đổi userID từ string sang ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("invalid user ID format: %v", err)
	}

	// Tìm kiếm người dùng theo ObjectID
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, nil // Không tìm thấy người dùng
		}
		return models.User{}, fmt.Errorf("error fetching user: %v", err)
	}

	// Trả về thông tin người dùng và không có lỗi
	return user, nil
}

// Lấy danh sách cảnh báo của người dùng từ MongoDB
func GetUserAlerts(userID string) ([]alert.Alert, error) {
	// Tạo bộ lọc để tìm các cảnh báo của người dùng
	filter := bson.M{"user_id": userID}

	// Tạo một slice để lưu các kết quả
	var alerts []alert.Alert

	// Truy vấn dữ liệu từ MongoDB
	cursor, err := config.AlertCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find alerts: %v", err)
	}
	defer cursor.Close(context.Background())

	// Lặp qua kết quả trả về và thêm vào slice alerts
	for cursor.Next(context.Background()) {
		var alert alert.Alert
		if err := cursor.Decode(&alert); err != nil {
			return nil, fmt.Errorf("failed to decode alert: %v", err)
		}
		alerts = append(alerts, alert)
	}

	// Kiểm tra lỗi từ việc lặp qua cursor
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return alerts, nil
}
