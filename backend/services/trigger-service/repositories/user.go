package repositories

import (
	"context"
	"fmt"
	"time"

	models "github.com/dath-241/coin-price-be-go/services/trigger-service/models/user"
	alert "github.com/dath-241/coin-price-be-go/services/trigger-service/models/alert"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByID(userID string) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tìm kiếm người dùng trong MongoDB theo ID
	err := utils.AlertCollection.FindOne(ctx, bson.M{"id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, nil // Trả về nil nếu không tìm thấy người dùng
		}
		return models.User{}, err
	}

	return user, nil
}

// Lấy danh sách cảnh báo của người dùng từ MongoDB
func GetUserAlerts(userID string) ([]alert.Alert, error) {
	// Tạo bộ lọc để tìm các cảnh báo của người dùng
	filter := bson.M{"user_id": userID}

	// Tạo một slice để lưu các kết quả
	var alerts []alert.Alert

	// Truy vấn dữ liệu từ MongoDB
	cursor, err := utils.AlertCollection.Find(context.Background(), filter)
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
