package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllUsers godoc
// @Summary Get all users
// @Description Admin can retrieve a list of all users in the system. Returns basic information about each user.
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 {array} models.UserResponse "List of users successfully fetched"
// @Failure 500 {object} models.ErrorResponse "Internal server error while fetching users"
// @Router /api/v1/admin/users [get]
func GetAllUsers(repo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Sử dụng repository để tìm dữ liệu
		users, err := repo.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error", // Failed to fetch users
			})
			return
		}

		// Tạo kết quả cần trả về
		var result []gin.H
		for _, user := range users {
			// Lấy thông tin người dùng và thêm vào kết quả
			result = append(result, gin.H{
				"user_id":  user.ID.Hex(),
				"username":	user.Username,
				"email":    user.Email,
				"vip_level": user.Role,
				"status":   user.IsActive, 
			})
		}

		// Trả về danh sách người dùng
		c.JSON(http.StatusOK, result)
	}
}

// GetUserByAdmin godoc
// @Summary Get user details by admin
// @Description Admin can retrieve user details by providing the user ID. Returns user information such as username, email, VIP level, etc.
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param user_id path string true "User ID"
// @Success 200 {object} models.UserDTO "User details successfully fetched"
// @Failure 400 {object} models.ErrorResponse "Invalid user ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error while fetching user"
// @Router /api/v1/admin/user/{user_id} [get]
func GetUserByAdmin(repo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id") // Lấy user_id từ URL
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID is required",
			})
			return
		}

		// Kiểm tra tính hợp lệ của ObjectID
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Tìm kiếm người dùng qua repository
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.UserDTO
		err = repo.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
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
			"user_id":    user.ID.Hex(),
			"username":   user.Username,
			"profile":    user.Profile,
			"email":      user.Email,
			"vip_level":  user.Role,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}
}

// DeleteUserByAdmin godoc
// @Summary Delete a user by admin
// @Description Admin can delete a user from the system by providing the user ID. The user will be permanently removed from the database.
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param user_id path string true "User ID"
// @Success 200 {object} models.MessageResponse "User deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid user ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error while deleting user"
// @Router /api/v1/admin/user/{user_id} [delete]
func DeleteUserByAdmin(repo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		userID := c.Param("user_id")

		// Kiểm tra tính hợp lệ của ObjectID
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Xóa người dùng
		result, err := repo.DeleteOne(context.Background(), bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error", // Failed to delete user
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

// GetPaymentHistoryForAdmin godoc
// @Summary Get payment history for all users (admin only)
// @Description Retrieves the payment history for all users. Admin can view all payments made by users.
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} models.PaymentHistoryAdminResponse "List of all payment histories"
// @Failure 500 {object} models.ErrorResponse "Internal server error while fetching payment history"
// @Router /api/v1/admin/payment-history [get]
func GetPaymentHistoryForAdmin(repo repository.PaymentRepository) func(*gin.Context) {
    return func(c *gin.Context) {
        // Lấy lịch sử thanh toán từ repository thay vì MongoDB trực tiếp
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        payments, err := repo.FindPayments(ctx, bson.M{})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal Server Error", // Error fetching payment history
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
                "order_id":          payment.OrderID,
                "user_id":           payment.UserID.Hex(),
                "orderInfo":         payment.OrderInfo,           // Thông tin đơn hàng
                "transaction_status": payment.TransactionStatus,  // Trạng thái giao dịch
                "amount":            payment.Amount,              // Số tiền thanh toán
            }
            paymentHistory = append(paymentHistory, paymentDetails)
        }

        c.JSON(http.StatusOK, gin.H{
            "payment_history": paymentHistory,
        })
    }
}

// GetPaymentHistoryForUserByAdmin godoc
// @Summary Retrieve payment history for a specific user (admin access only)
// @Description Retrieves the payment history of a specific user based on the provided user ID. This endpoint is restricted to admin access only and returns a list of payment transactions associated with the user.
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param user_id path string true "User ID" 
// @Success 200 {object} models.PaymentHisDetailsAdminResponse "List of payment transactions for the user"
// @Failure 400 {object} models.ErrorResponse "Invalid user ID or missing parameters"
// @Failure 404 {object} models.ErrorResponse "No payment history found for this user"
// @Failure 500 {object} models.ErrorResponse "Internal server error during payment history retrieval"
// @Router /api/v1/admin/payment-history/{user_id} [get]
func GetPaymentHistoryOfAUserByAdmin(repo repository.PaymentRepository) func(*gin.Context) {
    return func(c *gin.Context) {
        // Lấy user_id từ URL
        userID := c.Param("user_id")
        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "User ID is required",
            })
            return
        }

        // Lấy lịch sử thanh toán từ repository thay vì MongoDB trực tiếp
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        payments, err := repo.FindPayments(ctx, bson.M{"user_id": userID})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal Server Error", // Error fetching payment history
            })
            return
        }

        if len(payments) == 0 {
            c.JSON(http.StatusOK, gin.H{
                "message": "No payment history found for this user",
            })
            return
        }

        // Trả về danh sách các lịch sử thanh toán
        var paymentHistory []map[string]interface{}
        for _, payment := range payments {
            paymentDetails := map[string]interface{}{
                "order_id":          	payment.OrderID,
                //"user_id":           payment.UserID.Hex(),
                "orderInfo":        	payment.OrderInfo,           // Thông tin đơn hàng
                "transaction_status": 	payment.TransactionStatus,  // Trạng thái giao dịch
                "amount":            	payment.Amount,              // Số tiền thanh toán
				"vip_level":			payment.VipLevel,
				"payment_url":			payment.PaymentURL,
				"created_at":			payment.CreatedAt,	
				"updated_at":			payment.UpdatedAt,
            }
            paymentHistory = append(paymentHistory, paymentDetails)
        }

        c.JSON(http.StatusOK, gin.H{
            "payment_history": paymentHistory,
        })
    }
}

// BanAccount godoc
// @Summary Ban a user account
// @Description Ban a user account by setting the account status to inactive. Admin can use this to ban a user account.
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param user_id path string true "User ID"
// @Success 200 {object} models.MessageResponse "Account banned successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid user ID or missing user ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error while banning account"
// @Router /api/v1/admin/user/{user_id}/ban [put] 
func BanAccount(userRepo repository.UserRepository) func(*gin.Context) {
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

		// Chuyển đổi userID thành ObjectID
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Tìm người dùng trong cơ sở dữ liệu
		userResult := userRepo.FindOne(c, bson.M{"_id": objID})
		var user models.User
		if err := userResult.Decode(&user); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		// Cập nhật trạng thái tài khoản là đã bị cấm
		update := bson.M{
			"$set": bson.M{
				"is_active":  false,
				"updated_at": time.Now(),
			},
		}

		// Sử dụng repository để cập nhật trạng thái tài khoản
		updateResult, err := userRepo.UpdateOne(c, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error", // Thất bại khi cập nhật trạng thái tài khoản
			})
			return
		}

		if updateResult.MatchedCount == 0 {
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

// ActiveAccount godoc
// @Summary Activate a user account
// @Description Activate a user account by setting the account status to active. Admin can use this to activate a banned account.
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param user_id path string true "User ID"
// @Success 200 {object} models.MessageResponse "Account activated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid user ID or missing user ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error while activating account"
// @Router /api/v1/admin/user/{user_id}/active [put] 
func ActiveAccount(userRepo repository.UserRepository) func(*gin.Context) {
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

		// Chuyển đổi userID thành ObjectID
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Tìm người dùng trong cơ sở dữ liệu
		userResult := userRepo.FindOne(c, bson.M{"_id": objID})
		var user models.User
		if err := userResult.Decode(&user); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		// Cập nhật trạng thái tài khoản là đã bị cấm
		update := bson.M{
			"$set": bson.M{
				"is_active":  true,
				"updated_at": time.Now(),
			},
		}

		// Sử dụng repository để cập nhật trạng thái tài khoản
		updateResult, err := userRepo.UpdateOne(c, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error", // Thất bại khi cập nhật trạng thái tài khoản
			})
			return
		}

		if updateResult.MatchedCount == 0 {
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