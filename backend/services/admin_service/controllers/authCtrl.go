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

// Khởi tạo 1 user
func newUser(user models.User) models.User {
	// Gán giá trị mặc định
	user.Role = "VIP-0"                                        // Mặc định là 'VIP-0'
	user.IsActive = true                                       // Mặc định là true
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now()) // Mặc định thời gian hiện tại
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Khởi tạo Profile mặc định
	if user.Profile.FullName == "" {
		user.Profile.FullName = user.Username
	}
	if user.Profile.AvatarURL == "" {
		user.Profile.AvatarURL = "https://drive.google.com/file/d/15Ef4yebpGhT8pwgnt__utSESZtJdmA4a/view?usp=sharing"
	}
	return user 
}


// @Summary Register a new user
// @Description This endpoint allows a new user to register by providing a username, password, email, and optional phone number,...
// @Tags Authentication
// @Accept json
// @Produce json
// @Param RegisterRequest body models.RegisterRequest true "Register request body"
// @Success 201 {object}  models.MessageResponse "User registered successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request: Invalid input"
// @Failure 409 {object} models.ErrorResponse "Conflict: Email, username, or phone already exists"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error: Database operation failed"
// @Router /api/v1/auth/register [post]
func Register(userRepo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		var user models.User

		// Kiểm tra nhận được file JSON
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input",
			})
			return
		}

		// Kiểm tra định dạng username
		if !utils.IsValidUsername(user.Username) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Username only alphanumeric characters and hyphens are allowed.",
			})
			return
		}

		// Kiểm tra mật khẩu
		if !utils.IsValidPassword(user.Password) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password must contain at least 8 characters, including letters, numbers, and special characters.",
			})
			return
		}

		// Kiểm tra phone number nếu có
		if user.Profile.PhoneNumber != "" {
			if !utils.IsValidPhoneNumber(user.Profile.PhoneNumber) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid phone number.",
				})
				return
			}
		}

		// Kiểm tra xem username, email hoặc phone đã tồn tại chưa
		filter := bson.M{"$or": []bson.M{
			{"username": user.Username},
			{"email": user.Email},
		}}
		if user.Profile.PhoneNumber != "" {
			filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"profile.phone_number": user.Profile.PhoneNumber})
		}

		exists, err := userRepo.ExistsByFilter(context.Background(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email, username, or phone already exists.",
			})
			return
		}

		// Hash mật khẩu
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		user.Password = string(hashedPassword)

		// Tạo đối tượng user hoàn chỉnh
		user = newUser(user)

		// Thêm user vào database
		_, err = userRepo.InsertOne(context.Background(), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Response: 201 Created
		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
		})
	}
}

// Login godoc
// @Summary User Login
// @Description Authenticates a user by username or email and returns an token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "Login request body"
// @Success 200 {object} models.RorLResponse "Login successful"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Incorrect username or password"
// @Failure 403 {object} models.ErrorResponse "Forbidden: Account is inactive"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/auth/login [post]
func Login(userRepo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		var loginRequest models.LoginRequest

		// Kiểm tra input
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Tìm kiếm user bằng email hoặc username
		filter := bson.M{
			"$or": []bson.M{
				{"email": loginRequest.Username},
				{"username": loginRequest.Username},
			},
		}

		// Sử dụng repository để tìm user
		userResult := userRepo.FindOne(context.Background(), filter)
		var user models.User
		err := userResult.Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Username or password is incorrect",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
			}
			return
		}

		// Kiểm tra trạng thái is_active
		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Your account has been banned. Please contact support for assistance.",
			})
			return
		}

		// Kiểm tra mật khẩu
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Username or password is incorrect",
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

		// Trả về đăng nhập thành công và token
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
		})
	}
}

// Logout logs out the user by invalidating their token.
// @Summary Logout user
// @Description This API allows a user to log out by blacklisting their token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <JWT Token>"
// @Success 200 {object} models.MessageResponse "Logout successful"
// @Failure 400 {object} models.ErrorResponse "No token provided"
// @Router /api/v1/auth/logout [post]
// @Security Bearer
func Logout() func(*gin.Context) {
	return func(c *gin.Context) {
		// Lấy token từ header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
			return
		}

		// Xác thực Token
		accessClaims, err := middlewares.VerifyJWT(tokenString)
		if err == nil { // Token hợp lệ
			// Lấy thời gian hết hạn và thêm Token vào blacklist
			middlewares.BlacklistedTokens[tokenString] = accessClaims.ExpiresAt.Time
		}

		// Trả về thông báo thành công
		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
	}
}

// ForgotPassword send reset OTP to email.
// @Summary Request a password reset OTP
// @Description This API allows user can forgotPassword by sends a password reset OTP to the user's email address.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param ForgotPasswordRequest body models.ForgotPasswordRequest true "Forgot Password Request body"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse "Invalid request format"
// @Failure 404 {object} models.ErrorResponse "User not found with this email"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/auth/forgot-password [post]
func ForgotPassword(userRepo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		var request struct {
			Email string `json:"email" binding:"required,email"`
		}

		// Kiểm tra dữ liệu đầu vào
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Tìm user qua email
		filter := bson.M{"email": request.Email}
		userResult := userRepo.FindOne(context.Background(), filter)

		var user models.User
		err := userResult.Decode(&user)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found with this email",
			})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Tạo otp ngẫu nhiên
		otp, err := utils.GenerateOTP(6)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		hashedOTP := utils.HashString(otp)
		expiresAt := time.Now().Add(5 * time.Minute)

		// Cập nhật token vào database
		update := bson.M{
			"$set": bson.M{
				"reset_password_otp":   hashedOTP,
				"reset_password_expires": expiresAt,
			},
		}
		_, err = userRepo.UpdateOne(context.Background(), bson.M{"_id": user.ID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Gửi email
		if err := utils.SendEmail(request.Email, "Password Reset Request", user.Username, otp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Phản hồi thành công
		c.JSON(http.StatusOK, gin.H{"message": "Password reset OTP sent to your email"})
	}
}

// ResetPassword updates user's password using the provided OTP.
// @Summary Reset user password
// @Description This API allows users to reset their password using a valid OTP and a new password. The OTP is validated for authenticity and expiry before updating the password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param ResetRequest body models.ResetRequest true "Reset Password Request body"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse "Invalid request format or weak password"
// @Failure 401 {object} models.ErrorResponse "Invalid or expired OTP"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/auth/reset-password [post]
func ResetPassword(userRepo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		var request struct {
			OTP         string `json:"otp" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}

		// Kiểm tra dữ liệu đầu vào
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Kiểm tra độ mạnh mật khẩu
		if !utils.IsValidPassword(request.NewPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password must contain at least 8 characters, including letters, numbers, and special characters.",
			})
			return
		}

		// Hash OTP từ request
		hashedOTP := utils.HashString(request.OTP)

		// Tìm người dùng dựa trên OTP
		filter := bson.M{
			"reset_password_otp":   hashedOTP,
			"reset_password_expires": bson.M{"$gt": time.Now()}, // OTP còn hiệu lực
		}
		userResult := userRepo.FindOne(context.Background(), filter)

		var user models.User
		err := userResult.Decode(&user)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Hash mật khẩu mới
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Cập nhật mật khẩu và xóa OTP
		update := bson.M{
			"$set": bson.M{"password": hashedPassword},
			"$unset": bson.M{
				"reset_password_otp":   "",
				"reset_password_expires": "",
			},
		}
		_, err = userRepo.UpdateOne(context.Background(), bson.M{"_id": user.ID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Phản hồi thành công
		c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
	}
}


// RefreshToken generates a new token using the provided token.
// @Summary Refresh token
// @Description This API allows users to refresh their token using a valid old token. If the old token is valid and not blacklisted, a new token is generated.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <JWT Token>"
// @Security Bearer
// @Success 200 {object} models.RorLResponse "Token refreshed successfully"
// @Failure 401 {object} models.ErrorResponse "Token is missing, invalid, or blacklisted"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/auth/refresh-token [post]
func RefreshToken() func(*gin.Context) {
	return func(c *gin.Context) {
		// Lấy token từ header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Kiểm tra xem token có trong blacklist không
		if _, found := middlewares.BlacklistedTokens[tokenString]; found {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token has been blacklisted",
			})
			return
		}

		// Xác thực token và cấp lại token mới nếu hợp lệ
		tokenClaims, err := middlewares.VerifyJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		// Tạo mới token từ claims của token cũ
		newToken, err := middlewares.GenerateToken(tokenClaims.UserID, tokenClaims.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Đưa token cũ vào blacklist
		middlewares.BlacklistedTokens[tokenString] = tokenClaims.ExpiresAt.Time

		// Trả về thông báo thành công và token mới
		c.JSON(http.StatusOK, gin.H{
			"message": "Token refreshed successfully",
			"token":   newToken,
		})
	}
}
