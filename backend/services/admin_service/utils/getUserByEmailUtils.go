package utils

import (
    "fmt"
    "log"
	"context"

	"github.com/dath-241/coin-price-be-go/services/admin_service/config"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// Hàm lấy thông tin người dùng từ DB dựa trên email
func GetUserByEmail(email string) (*models.User, error) {
    // Lấy collection "users" từ DB
    collection := config.DB.Collection("User")

    var user models.User
    filter := bson.M{"email": email}

    // Truy vấn tìm kiếm người dùng theo email
    err := collection.FindOne(context.TODO(), filter).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            // Nếu không tìm thấy người dùng, trả về lỗi không có tài liệu
            return nil, fmt.Errorf("user not found with email: %s", email)
        }
        log.Println("Error retrieving user:", err)
        return nil, err
    }

    return &user, nil
}