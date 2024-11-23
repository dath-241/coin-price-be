package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/repository"
	"github.com/dath-241/coin-price-be-go/services/admin_service/utils"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)


// GetCurrentUserInfo retrieves the information of the currently authenticated user.
//
// @Summary Retrieve current user information
// @Description This endpoint fetches details of the currently authenticated user using the JWT token provided in the Authorization header.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <JWT Token>" 
// @Success 200 {object} models.UserDTO "Success: Returns the user's details"
// @Failure 400 {object} models.ErrorResponse "Bad Request: Invalid user ID"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: invalid token"
// @Failure 404 {object} models.ErrorResponse "Not Found: User not found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error: Failed to fetch user"
// @Router /api/v1/user/me [get]
// Lấy thông tin tài khoản người dùng hiện tại.
func GetCurrentUserInfo(repo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		// Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return

			// // Lỗi 401: Unauthorized
			// c.JSON(http.StatusUnauthorized, models.ErrorResponse{
    		// 	Error: "Authorization header required",
			// })
			// return 
		}

		// Xác thực token
		claims, err := middlewares.VerifyJWT(tokenString, true) // true indicates AccessToken
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
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

		// Chuyển user_id thành ObjectID
		objID, err := primitive.ObjectIDFromHex(currentUserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Tìm kiếm người dùng qua repository
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Tìm kiếm người dùng theo ObjectID using the repo interface
		filter := bson.M{"_id": objID}
		result := repo.FindOne(ctx, filter)
		if err := result.Err(); err != nil {
		//if err != nil {
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

		// Decode the result
		var user models.User
		if err := result.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error decoding user data",
			})
			return
		}

		// Trả về thông tin người dùng
		c.JSON(http.StatusOK, gin.H{
			"name":           user.Profile.FullName,
			"username":       user.Username,
			"email":          user.Email,
			"phone_number":   user.Profile.PhoneNumber,
			"avatar":         user.Profile.AvatarURL,
			"bio":            user.Profile.Bio,
			"date_of_birth":  user.Profile.DateOfBirth,
			"vip_level":      user.Role,
		})
	}
}


type UpdateUserProfileRequest struct {
    Name        string `json:"name" example:"John Doe"`
    Username    string `json:"username" example:"johndoe123"`
    Phone       string `json:"phone" example:"+1234567890"`
    Avatar      string `json:"avatar" example:"https://example.com/avatar.jpg"`
    Bio         string `json:"bio" example:"Software Engineer"`
    DateOfBirth string `json:"dateOfBirth" example:"2000-01-01"` // Định dạng: YYYY-MM-DD
}

type UpdateUserProfileResponse struct {
	 Message string `json:"message"`
}

// UpdateUserProfile updates the information of the currently authenticated user.
//
// @Summary Update current user profile
// @Description This endpoint allows the user to update their profile information such as name, username, phone, avatar, bio, and date of birth.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" Format("Bearer {token}")
// @Param UpdateUserProfileRequest body UpdateUserProfileRequest true "Update UserProfile Request body"
// @Success 200 {object} UpdateUserProfileResponse
// @Failure 400 {object} models.ErrorResponse "Bad Request: Invalid input"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Authorization header required or invalid token"
// @Failure 404 {object} models.ErrorResponse "Not Found: User not found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error: Failed to update user"
// @Router /api/v1/user/me [put]
// Chỉnh sửa thông tin tài khoản người dùng.
func UpdateUserProfile(repo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		var updateRequest struct {
			Name        string `json:"name"`
			Username    string `json:"username"`
			Phone       string `json:"phone"`
			Avatar      string `json:"avatar"`
			Bio         string `json:"bio"`
			DateOfBirth string `json:"dateOfBirth"` // Định dạng: YYYY-MM-DD
		}

		if err := c.ShouldBindJSON(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input",
			})
			return
		}

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		claims, err := middlewares.VerifyJWT(tokenString, true)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		currentUserID := claims.UserID
		if currentUserID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found in token",
			})
			return
		}

		objID, err := primitive.ObjectIDFromHex(currentUserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID format",
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
			filter["_id"] = bson.M{"$ne": objID} // Loại bỏ người dùng hiện tại
			users, err := repo.Find(c, filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error checking username or phone",
				})
				return
			}

			if len(users) > 0 {
				for _, user := range users {
					existingUser := user

					if updateRequest.Username != "" && updateRequest.Username == existingUser.Username {
						c.JSON(http.StatusConflict, gin.H{
							"error": "Username already in use",
						})
						return
					}

					if updateRequest.Phone != "" && updateRequest.Phone == existingUser.Profile.PhoneNumber {
						c.JSON(http.StatusConflict, gin.H{
							"error": "Phone number already in use",
						})
						return
					}
				}
			}
		}

		// Cập nhật chỉ các trường không trống
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

		updateResult, err := repo.UpdateOne(c, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update user",
			})
			return
		}

		if updateResult.MatchedCount == 0 {
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


type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}


// ChangePassword godoc
// @Summary Change user password
// @Description This endpoint allows an authenticated user to change their password by providing the current and new passwords.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" Format("Bearer {token}")
// @Param ChangePasswordRequest body ChangePasswordRequest true "Change Password Request body"
// @Success 200 {object} ChangePasswordResponse
// @Failure 400 {object} models.ErrorResponse "Bad Request: Password must contain at least 8 characters, including letters, numbers, and special characters."
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Authorization header required or invalid token"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Current password is incorrect"
// @Failure 404 {object} models.ErrorResponse "Not Found: User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/user/me/change_password [put]
// Đổi mật khẩu.
func ChangePassword(repo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		var request struct {
			CurrentPassword string `json:"current_password" binding:"required"`
			NewPassword     string `json:"new_password" binding:"required"`
		}

		// Parse JSON request
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Kiểm tra định dạng mật khẩu mới
		if !utils.IsValidPassword(request.NewPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password must contain at least 8 characters, including letters, numbers, and special characters.",
			})
			return
		}

		// Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Xác thực token
		claims, err := middlewares.VerifyJWT(tokenString, true)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Lấy userID từ token
		currentUserID := claims.UserID
		if currentUserID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		objID, err := primitive.ObjectIDFromHex(currentUserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Lấy thông tin người dùng từ repo
		// user, err := repo.FindOne(c, bson.M{"_id": objID})
		// if err != nil {
		// 	if err == mongo.ErrNoDocuments {
		// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		// 	} else {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		// 	}
		// 	return
		// }
		var user models.User
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

		// Kiểm tra mật khẩu hiện tại
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.CurrentPassword))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
			return
		}

		// Hash mật khẩu mới
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
			return
		}

		// Cập nhật mật khẩu
		update := bson.M{
			"$set": bson.M{
				"password":    string(hashedPassword),
				"updated_at": time.Now(),
			},
		}
		result, err := repo.UpdateOne(c, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Trả về kết quả thành công
		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	}
}


type ChangeMailRequest struct {
	Email string `json:"email" binding:"required,email"` // email mới cần cập nhật
}

type ChangeMailResponse struct {
	Message string `json:"message"`
}

// ChangeEmail godoc
// @Summary Change user email
// @Description This endpoint allows an authenticated user to change their email address.
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" Format("Bearer {token}")
// @Param ChangeMailRequest body ChangeMailRequest true "Change Mail Request body"
// @Success 200 {object} ChangeMailResponse
// @Failure 400 {object} models.ErrorResponse "Invalid email format"
// @Failure 401 {object} models.ErrorResponse "Unauthorized or token is invalid"
// @Failure 409 {object} models.ErrorResponse "Email already exists"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/user/me/change_email [put]
// Đổi email.
func ChangeEmail(userRepo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		var request struct {
			Email string `json:"email" binding:"required,email"` // email mới cần cập nhật
		}

		// Parse JSON để lấy dữ liệu cập nhật
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid email format",
			})
			return
		}

		// Lấy token từ header Authorization
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

		// Chuyển đổi userID thành ObjectID
		objID, err := primitive.ObjectIDFromHex(currentUserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Kiểm tra email đã tồn tại trong hệ thống (trừ chính người dùng hiện tại)
		existingUserResult := userRepo.FindOne(c, bson.M{
			"email": request.Email,
			"_id":   bson.M{"$ne": objID}, // Loại trừ người dùng hiện tại
		})
		
		var existingUser models.User
		err = existingUserResult.Decode(&existingUser)
		if err == nil {
			// Nếu email đã tồn tại
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email already exists.",
			})
			return
		}

		// Cập nhật email mới vào cơ sở dữ liệu
		update := bson.M{
			"$set": bson.M{
				"email":      request.Email,
				"updated_at": time.Now(),
			},
		}

		// Sử dụng repository để cập nhật thông tin
		updateResult, err := userRepo.UpdateOne(c, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update email",
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
			"message": "Email updated successfully",
		})
	}
}

type DeleteResponse struct {
	Message string `json:"message"`
}

// DeleteCurrentUser godoc
// @Summary Delete current user account
// @Description This endpoint allows the user to Delete the account of the currently authenticated user based on the access token.
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" Format("Bearer {token}")
// @Success 200 {object} DeleteResponse
// @Failure 400 {object} models.ErrorResponse "Invalid user ID"
// @Failure 401 {object} models.ErrorResponse "Unauthorized or token is invalid"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Failed to delete user"
// @Router /api/v1/user/me [delete]
// Xóa tài khoản người dùng hiện tại.
func DeleteCurrentUser(repo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		// Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
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

		// Chuyển đổi userID thành ObjectID
		objID, err := primitive.ObjectIDFromHex(currentUserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		// Xóa người dùng từ repository
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := repo.DeleteOne(ctx, bson.M{"_id": objID})
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



// GetPaymentHistory retrieves the payment history of the authenticated user.
//
// @Summary Retrieve payment history
// @Description This endpoint returns the payment history of the currently authenticated user using their JWT token.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <JWT Token>"
// @Success 200 {object} []map[string]interface{} "Success: Returns the user's payment history"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Invalid or missing token"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error: Failed to fetch payment history"
// @Router /api/v1/user/me/payment-history [get]
// Xem lịch sử thanh toán của người dùng
func GetPaymentHistory(repo repository.PaymentRepository) func(*gin.Context) {
    return func(c *gin.Context) {
        // Lấy token từ header Authorization
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

        // Tìm lịch sử thanh toán của người dùng theo userID từ repository
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        payments, err := repo.FindPayments(ctx, bson.M{"user_id": userID})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Error fetching payment history",
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
                "CreatedAt":        payment.CreatedAt,
                "UpdatedAt":        payment.UpdatedAt,
            }
            paymentHistory = append(paymentHistory, paymentDetails)
        }

        c.JSON(http.StatusOK, gin.H{
            "payment_history": paymentHistory,
        })
    }
}

