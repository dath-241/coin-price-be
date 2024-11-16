package controllers

import (
	"context"
	"net/http"
	"time"

	"backend/services/admin_service/src/config"
	"backend/services/admin_service/src/models"
	"backend/services/admin_service/src/utils"
	"backend/services/admin_service/src/middlewares"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Lấy thông tin tài khoản người dùng hiện tại.
func GetCurrentUserInfo() func(*gin.Context) {
	return func(c *gin.Context) {
		//userID := c.Param("user_id") // Lấy user_id từ URL

		// Lấy token từ header Authorization
		// tokenString := c.GetHeader("Authorization")
		// if tokenString == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		// 	return
		// }

		// Lấy token từ cookie
		tokenString, err := c.Cookie("accessToken")
		//fmt.Println("cookie", tokenString)
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token is required in cookies",
			})
			return
		}

		// Xác thực token
		claims, err := middlewares.VerifyJWT(tokenString, true) // true để chỉ định đây là AccessToken
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Lấy userID từ claims trong token
		currentUserID := claims.UserID
		if currentUserID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found in token",
			})
			return
		}

		// Kết nối đến database
		// if err := config.ConnectDatabase(); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": "Failed to connect to database",
		// 	})
		// 	return
		// }

		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User
		objID, err := primitive.ObjectIDFromHex(currentUserID) // Chuyển user_id thành ObjectID
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Tìm kiếm người dùng theo ObjectID
		err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to fetch user",
				})
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
		//var updatedData models.User
		var updateRequest struct {
			Name  string `json:"name"`
			Email string `json:"email" binding:"omitempty,email"`
		}
		// Parse JSON để lấy dữ liệu cập nhật
		if err := c.ShouldBindJSON(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if updateRequest.Email == "" && updateRequest.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid information",
			})
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

		// Lấy token từ header Authorization
		// tokenString := c.GetHeader("Authorization")
		// if tokenString == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		// 	return
		// }

		// Lấy token từ cookie
		tokenString, err := c.Cookie("accessToken")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token is required in cookies",
			})
			return
		}

		// Xác thực token
		claims, err := middlewares.VerifyJWT(tokenString, true) // true chỉ định đây là AccessToken
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Lấy userID từ claims
		currentUserID := claims.UserID
		if currentUserID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found in token",
			})
			return
		}

		// Kết nối đến database
		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		objID, err := primitive.ObjectIDFromHex(currentUserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Kiểm tra xem email hoặc tên có bị trùng với người dùng khác không
		filter := bson.M{}
		if updateRequest.Email != "" {
			filter["email"] = updateRequest.Email
		}
		if updateRequest.Name != "" {
			filter["name"] = updateRequest.Name
		}
		filter["_id"] = bson.M{"$ne": objID} // Đảm bảo không tìm thấy chính người dùng hiện tại

		// Kiểm tra nếu tên hoặc email đã tồn tại trong cơ sở dữ liệu
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error checking name or email",
			})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Name or email already in use",
			})
			return
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update user",
			})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
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
		// tokenString := c.GetHeader("Authorization")
		// if tokenString == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": "Authorization header required",
		// 	})
		// 	return
		// }

		// Lấy token từ cookie
		tokenString, err := c.Cookie("accessToken")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token is required in cookies",
			})
			return
		}

		// Xác thực token
		claims, err := middlewares.VerifyJWT(tokenString, true) // true chỉ định đây là AccessToken
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Lấy userID từ claims
		currentUserID := claims.UserID
		if currentUserID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found in token",
			})
			return
		}

		// Kết nối đến database
		// if err := config.ConnectDatabase(); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": "Failed to connect to database",
		// 	})
		// 	return
		// }

		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		objID, err := primitive.ObjectIDFromHex(currentUserID)
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
				"error": "Failed to delete user",
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
			"message": "User account deleted successfully",
		})
	}
}

// Xem lịch sử thanh toán của người dùng
func GetPaymentHistory() func(*gin.Context) {
    return func(c *gin.Context) {
        // Lấy token từ cookie
        tokenString, err := c.Cookie("accessToken")
        if err != nil || tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization token is required in cookies",
            })
            return
        }

        // Kiểm tra tính hợp lệ của token
        claims, err := middlewares.VerifyJWT(tokenString, true) // true để chỉ định đây là AccessToken
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": err.Error(),
            })
            return
        }

        // Lấy userID từ claims
        userID := claims.UserID
        if userID == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token claims",
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
                "error": "Error fetching payment history",
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

        // Trả về danh sách các lịch sử thanh toán
        c.JSON(http.StatusOK, gin.H{
            "payment_history": payments,
        })
    }
}
