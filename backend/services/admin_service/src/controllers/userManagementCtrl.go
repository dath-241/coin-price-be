package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"backend/services/admin_service/src/config"
	"backend/services/admin_service/src/models"
	"backend/services/admin_service/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Lấy thông tin tài khoản người dùng hiện tại.
func GetCurrentUserInfo() func(*gin.Context) {
	return func(c *gin.Context) {
		//userID := c.Param("user_id") // Lấy user_id từ URL

		// Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Kiểm tra tính hợp lệ của token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Lấy userID từ claims trong token
		currentUserID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		// Kết nối đến database
		if err := config.ConnectDatabase(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
			return
		}

		collection := config.DB.Collection("User")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User
		objID, err := primitive.ObjectIDFromHex(currentUserID) // Chuyển user_id thành ObjectID
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
			"email":      user.Email,
			"vip_level":  user.Role,
			"created_at": user.CreatedAt,
		})
	}
}

// Chỉnh sửa thông tin tài khoản người dùng.
func UpdateCurrentUser() func(*gin.Context) {
	return func(c *gin.Context) {
		// Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Kiểm tra tính hợp lệ của token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Lấy userID từ claims trong token
		currentUserID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		//var updatedData models.User
		var updateRequest struct {
			Name  string `json:"name"`
			Email string `json:"email" binding:"email"`
		}
		// Parse JSON để lấy dữ liệu cập nhật
		if err := c.ShouldBindJSON(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if updateRequest.Email == "" && updateRequest.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid information"})
			return
		}

		// Kiểm tra nếu có name thì có đúng định dạng chưa
		if updateRequest.Name != "" {
			if !utils.IsValidName(updateRequest.Name) {
				c.JSON(http.StatusBadRequest, gin.H{
				"error": "Name length must be between 1 and 50 characters.",
				})
				return
			}
			// Kiểm tra xem tên chỉ chứa các ký tự chữ cái
			if !utils.IsAlphabetical(updateRequest.Name) {
				c.JSON(http.StatusBadRequest, gin.H{
				"error": "Name must only contain alphabetical characters.",
				})
				return
			}
		}

		// Thêm check email đã tồn tại hay chưa

		// Kết nối đến database
		if err := config.ConnectDatabase(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
			return
		}

		collection := config.DB.Collection("User")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		objID, err := primitive.ObjectIDFromHex(currentUserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Kiểm tra email đã tồn tại hay chưa
		if updateRequest.Email != "" {
			filter := bson.M{"email": updateRequest.Email, "_id": bson.M{"$ne": objID}}
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking email"})
				return
			}
			if count > 0 {
				c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
				return
			}
		}

		// Tạo bản cập nhật chỉ với các trường không trống
		update := bson.M{"$set": bson.M{}}
		if updateRequest.Name != "" {
			update["$set"].(bson.M)["name"] = updateRequest.Name
		}
		if updateRequest.Email != "" {
			update["$set"].(bson.M)["email"] = updateRequest.Email
		}

		result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User information updated successfully",
		})
	}
}

// Xóa tài khoản người dùng hiện tại.
func DeleteCurrentUser() func(*gin.Context) {
	return func(c *gin.Context) {
		//userID := c.Param("user_id")
		// Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Kiểm tra tính hợp lệ của token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Lấy userID từ claims trong token
		currentUserID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		// Kết nối đến database
		if err := config.ConnectDatabase(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
			return
		}

		collection := config.DB.Collection("User")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		objID, err := primitive.ObjectIDFromHex(currentUserID)
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

		c.JSON(http.StatusOK, gin.H{"message": "User account deleted successfully"})
	}
}

// func UpgradeUserVIP() func(*gin.Context) {
//     return func(c *gin.Context) {
//         userID := c.Param("user_id")
//         var upgradeRequest struct {
//             PaymentID   string `json:"payment_id" binding:"required"`
//             NewVIPLevel string `json:"new_vip_level" binding:"required,oneof=VIP-1 VIP-2 VIP-3"`
//         }

//         // Parse JSON để lấy dữ liệu yêu cầu nâng cấp
//         if err := c.ShouldBindJSON(&upgradeRequest); err != nil {
//             c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//             return
//         }

//         // Kết nối đến database
//         if err := config.ConnectDatabase(); err != nil {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
//             return
//         }

//         collection := config.DB.Collection("User")

//         ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//         defer cancel()

//         objID, err := primitive.ObjectIDFromHex(userID)
//         if err != nil {
//             c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
//             return
//         }

//         // Cập nhật cấp độ VIP của người dùng
//         update := bson.M{
//             "$set": bson.M{
//                 "role":       upgradeRequest.NewVIPLevel,
//                 "upgraded_at": time.Now(),
//                 "payment_id": upgradeRequest.PaymentID,
//             },
//         }

//         opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
//         var updatedUser models.User
//         err = collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updatedUser)
//         if err != nil {
//             if err == mongo.ErrNoDocuments {
//                 c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//             } else {
//                 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade user VIP level"})
//             }
//             return
//         }

//         // Trả về thông tin người dùng sau khi nâng cấp VIP
//         c.JSON(http.StatusOK, gin.H{
//             "user_id":      userID,
//             "name":     updatedUser.Name,
//             "email":        updatedUser.Email,
//             "role":         updatedUser.Role,
//             "upgraded_at":  time.Now(),
//             "payment_id":   upgradeRequest.PaymentID,
//         })
//     }
// }
