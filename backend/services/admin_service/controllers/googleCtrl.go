package controllers

import (
	"context"
	"net/http"

	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/repository"
	"github.com/dath-241/coin-price-be-go/services/admin_service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// utlis
func generateUniqueUsername() string {
	return "user-" + uuid.New().String()
}


// GoogleLogin authenticates the user using their Google ID Token.
// @Summary Google Login
// @Description This endpoint allows users to authenticate using their Google account. The frontend sends a Google ID Token, which is verified on the backend to create or authenticate the user.
// @Tags Authentication
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id_token formData string true "Google ID Token"
// @Success 200 {object} models.RorLResponse "Success: Login successful with token"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Invalid Google ID token"
// @Failure 403 {object} models.ErrorResponse "Forbidden: User account is banned"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error: Failed to create or retrieve user"
// @Router /api/v1/auth/google-login [post]
// GoogleLogin xử lý đăng nhập bằng Google ID Token
func GoogleLogin(userRepo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		idToken := c.PostForm("id_token") // Nhận Google ID Token từ frontend

		// Xác minh Google ID Token
		userInfo, err := utils.VerifyGoogleIDToken(idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Google ID token",
			})
			return
		}

		// Lấy email từ thông tin Google
		email := userInfo["email"].(string)

		// Tìm người dùng trong cơ sở dữ liệu dựa trên email
		filter := bson.M{"email": email}
		userResult := userRepo.FindOne(context.Background(), filter)

		var user models.User
		err = userResult.Decode(&user)
		if err == mongo.ErrNoDocuments {
			// Nếu không tìm thấy user, tạo user mới
			user.Profile.FullName = userInfo["name"].(string)
			user.Profile.AvatarURL = userInfo["picture"].(string)
            user.Username = generateUniqueUsername()
            user.Email = email

            user = newUser(user)
			
			// Thêm user vào database
			insertResult, insertErr := userRepo.InsertOne(context.Background(), user)
			if insertErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
				return
			}
            user.ID = insertResult.InsertedID.(primitive.ObjectID)
		} else if err != nil {
			// Lỗi khi truy xuất database
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

        // Kiểm tra tài khoản có bị cấm hay không
		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Your account is banned",
			})
			return
		}

		// Tạo JWT Access token 
		token, err := middlewares.GenerateToken(user.ID.Hex(), user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Phản hồi thành công
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
		})
	}
}
