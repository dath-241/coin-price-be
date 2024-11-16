package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/dath-241/coin-price-be-go/services/admin_service/config"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Lấy thông tin tất cả người dùng bởi admin
func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Kết nối đến database
		// if err := config.ConnectDatabase(); err != nil {
		//     c.JSON(http.StatusInternalServerError, gin.H{
		//         "error": "Failed to connect to database",
		//     })
		//     return
		// }

		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Lấy tất cả người dùng
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch users",
			})
			return
		}
		defer cursor.Close(ctx)

		var users []models.User
		if err := cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to parse users",
			})
			return
		}

		// Trả về danh sách người dùng
		var result []gin.H
		for _, user := range users {
			result = append(result, gin.H{
				"user_id": user.ID.Hex(),
				//"username":   user.Name,
				"email":     user.Email,
				"vip_level": user.Role,
				"status":    "active",
			})
		}

		c.JSON(http.StatusOK, result)
	}
}

// Lấy thông tin của 1 người dùng bởi admin
func GetUserByAdmin() func(*gin.Context) {
	return func(c *gin.Context) {
		userID := c.Param("user_id") // Lấy user_id từ URL

		// Kết nối đến database
		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User
		objID, err := primitive.ObjectIDFromHex(userID) // Chuyển user_id thành ObjectID
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Tìm kiếm người dùng theo ObjectID
		err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			}
			return
		}

		// Trả về dữ liệu người dùng
		c.JSON(http.StatusOK, gin.H{
			"user_id":    user.ID,
			"username":   user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		})
	}
}

// Xóa 1 người dùng từ user_id bởi admin
func DeleteUserByAdmin() func(*gin.Context) {
	return func(c *gin.Context) {
		userID := c.Param("user_id")

		// Kết nối đến database
		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Xóa người dùng
		result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

// Xem lịch sử thanh toán của tất cả người dùng (dành cho admin)
func GetPaymentHistoryForAdmin() func(*gin.Context) {
	return func(c *gin.Context) {
		// Lấy lịch sử thanh toán từ MongoDB
		collection := config.DB.Collection("OrderMoMo")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error fetching payment history",
			})
			return
		}
		defer cursor.Close(ctx)

		var payments []models.Order
		for cursor.Next(ctx) {
			var payment models.Order
			if err := cursor.Decode(&payment); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error decoding payment data",
				})
				return
			}
			payments = append(payments, payment)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error reading from cursor",
			})
			return
		}

		if len(payments) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"message": "No payment history found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"payment_history": payments,
		})
	}
}
