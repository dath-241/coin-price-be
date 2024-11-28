package controllers

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
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
	return user // Bạn có thể thay đổi giá trị trả về dựa trên logic của bạn
}

// setAuthCookies sẽ thiết lập các cookie cho accessToken và refreshToken nếu được yêu cầu
func setAuthCookies(c *gin.Context, accessToken, refreshToken string, setAccessToken, setRefreshToken bool) error {
	// Load biến môi trường cho tên miền cookie và thời gian sống
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	accessTokenTTL := os.Getenv("ACCESS_TOKEN_TTL") // Thời gian sống token
	refreshTokenTTL := os.Getenv("REFRESH_TOKEN_TTL")

	if cookieDomain == "" || accessTokenTTL == "" || refreshTokenTTL == "" {
		return fmt.Errorf("environment variables are not set")
	}

	accessTokenTTLInt, err := strconv.Atoi(accessTokenTTL)
	if err != nil {
		return fmt.Errorf("invalid ACCESS_TOKEN_TTL format")
	}

	refreshTokenTTLInt, err := strconv.Atoi(refreshTokenTTL)
	if err != nil {
		return fmt.Errorf("invalid REFRESH_TOKEN_TTL format")
	}

	// Nếu set accessToken là true thì thiết lập cookie accessToken
	if setAccessToken {
		c.SetCookie("accessToken", accessToken, accessTokenTTLInt, "/api/v1", cookieDomain, true, true)      // chỉ dành cho /api/v1
		//c.SetCookie("accessToken", accessToken, accessTokenTTLInt, "/api/v1/auth/logout", cookieDomain, true, true) // chỉ dành cho /auth/logout
	}

	// Nếu set refreshToken là true thì thiết lập cookie refreshToken
	if setRefreshToken {
		c.SetCookie("refreshToken", refreshToken, refreshTokenTTLInt, "/api/v1/auth/refresh-token", cookieDomain, true, true)     // chỉ dành cho /auth/refresh-token
		c.SetCookie("refreshToken", refreshToken, refreshTokenTTLInt, "/api/v1/auth/logout", cookieDomain, true, true)            // chỉ dành cho /auth/logout
		//c.SetCookie("refreshToken", refreshToken, refreshTokenTTLInt, "/api/v1/payment/confirm", cookieDomain, true, true) // dành cho /api/v1/payment/confirm
	}

	return nil
}

// resetAuthCookies sẽ xóa các cookie cho accessToken và refreshToken
func resetAuthCookies(c *gin.Context) error {
	// Load biến môi trường cho tên miền cookie và thời gian sống
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDomain == "" {
		return fmt.Errorf("environment variables are not set")
	}

	// Xóa cookie Access Token và Refresh Token
	c.SetCookie("accessToken", "", 0, "/", cookieDomain, true, true)
	c.SetCookie("refreshToken", "", 0, "/", cookieDomain, true, true)

	return nil
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
				"error": "Error checking user existence",
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
				"error": "Failed to hash password",
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
				"error": "Failed to create user",
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
// @Description Authenticates a user by username or email and returns an access token
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
					"error": "Failed to find user",
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

		// Tạo JWT Refresh token
		refreshToken, err := middlewares.GenerateRefreshToken(user.ID.Hex(), user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate refresh token",
			})
			return
		}

		// Tạo JWT Access token
		accessToken, err := middlewares.GenerateAccessToken(user.ID.Hex(), user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate access token",
			})
			return
		}

		// Gọi hàm set cookie để thiết lập cookies cho người dùng
		err = setAuthCookies(c, accessToken, refreshToken, false, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Trả về đăng nhập thành công và token
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   accessToken,
		})
	}
}

// Logout logs out the user by invalidating their access and refresh tokens.
// @Summary Logout user
// @Description This API allows a user to log out by blacklisting their access and refresh tokens and clearing their authentication cookies.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <JWT Token>" 
// @Success 200 {object} models.MessageResponse "Logout successful"
// @Failure 400 {object} models.ErrorResponse "No token provided or Refresh Token not provided"
// @Failure 401 {object} models.ErrorResponse "Token has been revoked"
// @Failure 500 {object} models.ErrorResponse "Failed to reset cookies or other server error"
// @Router /api/v1/auth/logout [post]
// @Security Bearer
func Logout() func(*gin.Context) {
	return func(c *gin.Context) {
		// Lấy token từ header
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
			return
		}

		// // Lấy Access Token từ cookie
		// accessToken, err := c.Cookie("accessToken")
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": "Access Token not provided",
		// 	})
		// 	return
		// }

		// Lấy Refresh Token từ cookie
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Refresh Token not provided",
			})
			return
		}

		// Xác thực Access Token
		accessClaims, err := middlewares.VerifyJWT(accessToken, true)
		if err == nil { // Token hợp lệ
			// Lấy thời gian hết hạn và thêm Access Token vào blacklist
			middlewares.BlacklistedTokens[accessToken] = accessClaims.ExpiresAt.Time
		}

		// Xác thực Refresh Token
		refreshClaims, err := middlewares.VerifyJWT(refreshToken, false)
		if err == nil { // Token hợp lệ
			// Lấy thời gian hết hạn và thêm Refresh Token vào blacklist
			middlewares.BlacklistedTokens[refreshToken] = refreshClaims.ExpiresAt.Time
		}

		// Gọi hàm reset cookie để xóa cookies
		err = resetAuthCookies(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Trả về thông báo thành công
		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
	}
}

// ForgotPassword send reset link to email.
// @Summary Request a password reset link
// @Description This API allows user can forgotPassword by sends a password reset link to the user's email address.
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
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found with this email"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Tạo token ngẫu nhiên
		rawToken, err := utils.GenerateRandomString(32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
			return
		}
		hashedToken := utils.HashString(rawToken)
		expiresAt := time.Now().Add(15 * time.Minute)

		// Cập nhật token vào database
		update := bson.M{
			"$set": bson.M{
				"reset_password_token":   hashedToken,
				"reset_password_expires": expiresAt,
			},
		}
		_, err = userRepo.UpdateOne(context.Background(), bson.M{"_id": user.ID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset token"})
			return
		}

		// Tạo link reset password
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Base URL for reset password is missing"})
			return
		}
		resetLink := fmt.Sprintf("%s/reset-password?token=%s", baseURL, rawToken)

		// Chuẩn bị email template
		emailTemplatePath := "services/admin_service/templates/password_reset_email.html"
		htmlBody, err := os.ReadFile(emailTemplatePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read email template"})
			return
		}

		t, err := template.New("reset-email").Parse(string(htmlBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse email template"})
			return
		}

		var bodyBuffer bytes.Buffer
		err = t.Execute(&bodyBuffer, map[string]interface{}{
			"Name":      user.Username,
			"ResetLink": resetLink,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute email template"})
			return
		}

		// Gửi email
		htmlBodyString := bodyBuffer.String()
		if err := utils.SendEmail(request.Email, "Password Reset Request", htmlBodyString); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
			return
		}

		// Phản hồi thành công
		c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent to your email"})
	}
}

// ResetPassword updates user's password using the provided token.
// @Summary Reset user password
// @Description This API allows users to reset their password using a valid reset token and a new password. The token is validated for authenticity and expiry before updating the password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param ResetRequest body models.ResetRequest true "Reset Password Request body"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse "Invalid request format or weak password"
// @Failure 401 {object} models.ErrorResponse "Invalid or expired reset token"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/auth/reset-password [post]
func ResetPassword(userRepo repository.UserRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		var request struct {
			Token       string `json:"token" binding:"required"`
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

		// Hash token từ request
		hashedToken := utils.HashString(request.Token)

		// Tìm người dùng dựa trên token reset
		filter := bson.M{
			"reset_password_token":   hashedToken,
			"reset_password_expires": bson.M{"$gt": time.Now()}, // Token còn hiệu lực
		}
		userResult := userRepo.FindOne(context.Background(), filter)

		var user models.User
		err := userResult.Decode(&user)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired reset token"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Hash mật khẩu mới
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Cập nhật mật khẩu và xóa token reset
		update := bson.M{
			"$set": bson.M{"password": hashedPassword},
			"$unset": bson.M{
				"reset_password_token":   "",
				"reset_password_expires": "",
			},
		}
		_, err = userRepo.UpdateOne(context.Background(), bson.M{"_id": user.ID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		// Phản hồi thành công
		c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
	}
}

// RefreshToken generates a new access token using the provided refresh token.
// @Summary Refresh access token
// @Description This API allows users to refresh their access token using a valid refresh token stored in cookies. If the refresh token is valid and not blacklisted, a new access token is generated.
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} models.RorLResponse "New access token generated successfully"
// @Failure 401 {object} models.ErrorResponse "Refresh token is missing, invalid, or blacklisted"
// @Failure 500 {object} models.ErrorResponse "Internal server error while generating a new access token"
// @Router /api/v1/auth/refresh-token [post]
func RefreshToken() func(*gin.Context) {
	return func(c *gin.Context) {
		// Lấy token từ header Authorization
		// tokenString := c.GetHeader("Authorization")
		// if tokenString == "" {
		//     c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		//     return
		// }

		// Lấy refresh token từ cookie
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil || refreshToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Refresh token is missing or invalid",
			})
			return
		}

		// Kiểm tra xem Refresh Token có trong blacklist không
		if _, found := middlewares.BlacklistedTokens[refreshToken]; found {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Refresh token has been blacklisted",
			})
			return
		}

		// Xác thực Refresh Token và cấp lại Access Token nếu hợp lệ
		refreshClaims, err := middlewares.VerifyJWT(refreshToken, false)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid refresh token",
			})
			return
		}

		// Tạo mới Access Token từ claims của Refresh Token
		accessToken, err := middlewares.GenerateAccessToken(refreshClaims.UserID, refreshClaims.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate access token",
			})
			return
		}

		// Gọi hàm set cookie để thiết lập cookies cho người dùng
		// err = setAuthCookies(c, accessToken, refreshToken, true, false)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		// Trả về thông báo thành công và token mới
		c.JSON(http.StatusOK, gin.H{
			"message": "Access token refreshed successfully",
			"token":   accessToken,
		})
	}
}
