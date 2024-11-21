package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/dath-241/coin-price-be-go/services/admin_service/config"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/utils"
	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Lấy thông tin tài khoản người dùng hiện tại.
func GetCurrentUserInfo() func(*gin.Context) {
	return func(c *gin.Context) {
		//userID := c.Param("user_id") // Lấy user_id từ URL

		//Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Lấy token từ cookie
		// tokenString, err := c.Cookie("accessToken")
		// //fmt.Println("cookie", tokenString)
		// if err != nil || tokenString == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": "Authorization token is required in cookies",
		// 	})
		// 	return
		// }

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
			//"user_id":    user.ID,
			"name":			user.Profile.FullName,
			"username": 	user.Username,
			"email":      	user.Email,
			"phone_number": user.Profile.PhoneNumber,
			"avatar":		user.Profile.AvatarURL,
			"bio":			user.Profile.Bio,
			"date_of_birth":user.Profile.DateOfBirth,
			"vip_level":  	user.Role,
		})
	}
}

// Chỉnh sửa thông tin tài khoản người dùng.
func UpdateUserProfile() func(*gin.Context) {
	return func(c *gin.Context) {
		// Cấu trúc dữ liệu yêu cầu cập nhật
		var updateRequest struct {
			Name        string `json:"name"`
			Username    string `json:"username"`
			Phone       string `json:"phone"`
			Avatar      string `json:"avatar"`
			Bio         string `json:"bio"`
			DateOfBirth string `json:"dateOfBirth"` // Định dạng: YYYY-MM-DD
		}

		// Parse JSON để lấy dữ liệu cập nhật
		if err := c.ShouldBindJSON(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Lấy token từ cookie
		// tokenString, err := c.Cookie("accessToken")
		// if err != nil || tokenString == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": "Authorization token is required in cookies",
		// 	})
		// 	return
		// }

		//Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
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

		// Kiểm tra trùng lặp username hoặc phone
		filter := bson.M{"$or": []bson.M{}}
		if updateRequest.Username != "" {
			filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"username": updateRequest.Username})
		}
		if updateRequest.Phone != "" {
			filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"profile.phone_number": updateRequest.Phone})
		}
		if len(filter["$or"].([]bson.M)) > 0 {
			filter["_id"] = bson.M{"$ne": objID} // Không kiểm tra chính người dùng hiện tại
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error checking username or phone",
				})
				return
			}
			if count > 0 {
				if updateRequest.Username != "" {
					c.JSON(http.StatusConflict, gin.H{
						"error": "Username already in use",
					})
					return
				}
				if updateRequest.Phone != "" {
					c.JSON(http.StatusConflict, gin.H{
						"error": "Phone already in use",
					})
					return
				}
			}
		}

		// Tạo bản cập nhật chỉ với các trường không trống
		update := bson.M{
			"$set": bson.M{
				"updated_at": time.Now(), 
			},
		}
		if updateRequest.Name != "" {
			update["$set"].(bson.M)["profile.full_name"] = updateRequest.Name
		}
		if updateRequest.Username != "" {
			update["$set"].(bson.M)["username"] = updateRequest.Username
		}
		if updateRequest.Phone != "" {
			update["$set"].(bson.M)["profile.phone_number"] = updateRequest.Phone
		}
		if updateRequest.Avatar != "" {
			update["$set"].(bson.M)["profile.avatar_url"] = updateRequest.Avatar
		}
		if updateRequest.Bio != "" {
			update["$set"].(bson.M)["profile.bio"] = updateRequest.Bio
		}
		if updateRequest.DateOfBirth != "" {
			update["$set"].(bson.M)["profile.date_of_birth"] = updateRequest.DateOfBirth
		}

		// Cập nhật thông tin người dùng
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

// Đổi mật khẩu.
func ChangePassword() func(*gin.Context) {
	return func(c *gin.Context) {
		//var updatedData models.User
		var request struct {
			CurrentPassword string 	`json:"current_password" binding:"required"`
			NewPassword 	string 	`json:"new_password" binding:"required"`
		}
		// Parse JSON để lấy dữ liệu cập nhật
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// if request.CurrentPassword == "" || request.NewPassword == "" {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": "Invalid information",
		// 	})
		// 	return
		// }

		// Kiểm tra định dạng password mới
		if !utils.IsValidPassword(request.NewPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password must contain at least 8 characters, including letters, numbers, and special characters.",
			})
			return
		}

		// // Lấy token từ cookie
		// tokenString, err := c.Cookie("accessToken")
		// if err != nil || tokenString == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": "Authorization token is required in cookies",
		// 	})
		// 	return
		// }
		//Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
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

		// Lấy thông tin người dùng hiện tại từ database
		var user struct {
			Password string `bson:"password"`
		}
		err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		// Kiểm tra mật khẩu hiện tại
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.CurrentPassword))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Current password is incorrect",
			})
			return
		}

		// Hash mật khẩu mới 
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash new password",
			})
			return
		}

		// Lưu mk mới 
		update := bson.M{
			"$set": bson.M{
				"password": string(hashedPassword),
				"updated_at": time.Now(),
			},
		}

		result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update password",
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
			"message": "Password updated successfully",
		})
	}
}

// Đổi email.
func ChangeEmail() func(*gin.Context) {
	return func(c *gin.Context) {
		var request struct {
			Email string `json:"email" binding:"required,email"` // email mới cần cập nhật
		}

		// Parse JSON để lấy dữ liệu cập nhật
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Email không hợp lệ",
			})
			return
		}

		// // Lấy token từ cookie
		// tokenString, err := c.Cookie("accessToken")
		// if err != nil || tokenString == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": "Authorization token is required in cookies",
		// 	})
		// 	return
		// }

		//Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
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

		// Kiểm tra email đã tồn tại trong hệ thống (trừ chính người dùng hiện tại)
		var existingUser models.User
		err = collection.FindOne(ctx, bson.M{
			"email": request.Email,
			"_id": bson.M{"$ne": objID}, // Loại trừ người dùng hiện tại
		}).Decode(&existingUser)

		if err == nil {
			// Nếu email đã tồn tại
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email already exists.",
			})
			return
		}

		// Lưu email mới vào cơ sở dữ liệu
		update := bson.M{
			"$set": bson.M{
				"email": request.Email,
				"updated_at": time.Now(),
			},
		}

		result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update email",
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
			"message": "Email updated successfully",
		})
	}
}

// Xóa tài khoản người dùng hiện tại.
func DeleteCurrentUser() func(*gin.Context) {
	return func(c *gin.Context) {
		//userID := c.Param("user_id")
		//Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			return
		}

		// Lấy token từ cookie
		// tokenString, err := c.Cookie("accessToken")
		// if err != nil || tokenString == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": "Authorization token is required in cookies",
		// 	})
		// 	return
		// }

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
        // tokenString, err := c.Cookie("accessToken")
        // if err != nil || tokenString == "" {
        //     c.JSON(http.StatusUnauthorized, gin.H{
        //         "error": "Authorization token is required in cookies",
        //     })
        //     return
        // }

		//Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
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
        var paymentHistory []map[string]interface{}
        for _, payment := range payments {
            paymentDetails := map[string]interface{}{
                "OrderInfo":        payment.OrderInfo,         // Thông tin đơn hàng
                "TransactionStatus": payment.TransactionStatus, // Trạng thái giao dịch
                "Amount":           payment.Amount,             // Số tiền thanh toán
				"CreatedAt":		payment.CreatedAt,
				"UpdateAt":			payment.UpdatedAt,
            }
            paymentHistory = append(paymentHistory, paymentDetails)
        }

        c.JSON(http.StatusOK, gin.H{
            "payment_history": paymentHistory,
        })
    }
}

