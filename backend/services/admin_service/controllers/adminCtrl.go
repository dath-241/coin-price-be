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
		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Lấy tất cả người dùng
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error", //Failed to fetch users
			})
			return
		}
		defer cursor.Close(ctx)

		var users []models.UserDTO
		if err := cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error", //Failed to parse users
			})
			return
		}

		// Trả về danh sách người dùng
		var result []gin.H
		for _, user := range users {
			result = append(result, gin.H{
				"user_id": user.ID.Hex(),
				"username":   user.Username,
				"email":     user.Email,
				"vip_level": user.Role,
				"status":    user.IsActive,
			})
		}

		c.JSON(http.StatusOK, result)
	}
}

// Lấy thông tin của 1 người dùng bởi admin
func GetUserByAdmin() func(*gin.Context) {
	return func(c *gin.Context) {
		userID := c.Param("user_id") // Lấy user_id từ URL
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID is required",
			})
			return
		}
		// Kiểm tra tính hợp lệ của ObjectID
		_, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Kết nối đến database
		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.UserDTO
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
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error", //Failed to fetch user
				})
			}
			return
		}

		// Trả về dữ liệu người dùng
		c.JSON(http.StatusOK, gin.H{
			"user_id": 		user.ID.Hex(),
			"username": 	user.Username,
			"profile": 		user.Profile,
			"email":    	user.Email,
			"vip_level":	user.Role,
			"is_active":   	user.IsActive,
			"created_at":	user.CreatedAt,
			"updated_at": 	user.UpdatedAt,
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Xóa người dùng
		result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error", //Failed to delete user
			})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User deleted successfully",
		})
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
				"error": "Internal Server Error", //Error fetching payment history
			})
			return
		}
		defer cursor.Close(ctx)

		var payments []models.Order
		for cursor.Next(ctx) {
			var payment models.Order
			if err := cursor.Decode(&payment); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error", //Error decoding payment data
				})
				return
			}
			payments = append(payments, payment)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error", //Error reading from cursor
			})
			return
		}

		if len(payments) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"message": "No payment history found",
			})
			return
		}

		// Trả về danh sách các lịch sử thanh toán
        var paymentHistory []map[string]interface{}
        for _, payment := range payments {
            paymentDetails := map[string]interface{}{
				"order_id":				payment.OrderID,	
				"user_id":				payment.UserID,	
                "orderInfo":        	payment.OrderInfo,         // Thông tin đơn hàng
                "transaction_status": 	payment.TransactionStatus, // Trạng thái giao dịch
                "amount":           	payment.Amount,             // Số tiền thanh toán
            }
            paymentHistory = append(paymentHistory, paymentDetails)
        }

        c.JSON(http.StatusOK, gin.H{
            "payment_history": paymentHistory,
        })
	}
}

// Xem lịch sử thanh toán của một người dùng (dành cho admin)
func GetPaymentHistoryForUserByAdmin() func(*gin.Context) {
	return func(c *gin.Context) {
		// Lấy user_id từ URL
		userID := c.Param("user_id")

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID is required",
			})
			return
		}
		// Kiểm tra tính hợp lệ của ObjectID
		_, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Kết nối tới cơ sở dữ liệu MongoDB
		collection := config.DB.Collection("OrderMoMo")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Tìm lịch sử thanh toán của người dùng theo userID
		cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error", //Error fetching payment history
			})
			return
		}
		defer cursor.Close(ctx)

		// Lưu các đơn hàng thanh toán
		var payments []models.Order
		for cursor.Next(ctx) {
			var payment models.Order
			if err := cursor.Decode(&payment); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error", //Error decoding payment data
				})
				return
			}
			payments = append(payments, payment)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error", //Error reading from cursor
			})
			return
		}

		// Nếu không có lịch sử thanh toán nào
		if len(payments) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"message": "No payment history found for this user",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"payment_history": payments,
		})
	}
}

// Ban account 
func BanAccount() func(*gin.Context) {
	return func(c *gin.Context) {
		userID := c.Param("user_id") // Lấy user_id từ URL
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID is required",
			})
			return
		}
		// Kiểm tra tính hợp lệ của ObjectID
		_, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Kết nối đến database
		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Tìm kiếm người dùng theo ObjectID
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Lưu email mới vào cơ sở dữ liệu
		update := bson.M{
			"$set": bson.M{
				"is_active": false,
				"updated_at": time.Now(),
			},
		}

		result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error", //Failed to update account status
			})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		// Trả về kết quả thành công
		c.JSON(http.StatusOK, gin.H{
			"message": "Account has been banned successfully",
		})
	}
}

// Ban account 
func ActiveAccount() func(*gin.Context) {
	return func(c *gin.Context) {
		userID := c.Param("user_id") // Lấy user_id từ URL
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID is required",
			})
			return
		}
		// Kiểm tra tính hợp lệ của ObjectID
		_, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Kết nối đến database
		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Tìm kiếm người dùng theo ObjectID
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Lưu email mới vào cơ sở dữ liệu
		update := bson.M{
			"$set": bson.M{
				"is_active": true,
				"updated_at": time.Now(),
			},
		}

		result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		// Trả về kết quả thành công
		c.JSON(http.StatusOK, gin.H{
			"message": "Account has been active successfully",
		})
	}
}