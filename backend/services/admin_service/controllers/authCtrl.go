package controllers

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dath-241/coin-price-be-go/services/admin_service/config"
	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/utils"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Khởi tạo 1 user 
func newUser(user models.User) models.User {
	// Gán giá trị mặc định
	user.Role = "VIP-0"       // Mặc định là 'VIP-0'
	user.IsActive = true        // Mặc định là true
	user.CreatedAt = time.Now() // Mặc định thời gian hiện tại
	user.UpdatedAt = time.Now()

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
		c.SetCookie("accessToken", accessToken, accessTokenTTLInt, "/api/v1", cookieDomain, true, true)                // chỉ dành cho /api/v1
		c.SetCookie("accessToken", accessToken, accessTokenTTLInt, "/auth/logout", cookieDomain, true, true)               // chỉ dành cho /auth/logout
	}

	// Nếu set refreshToken là true thì thiết lập cookie refreshToken
	if setRefreshToken {
		c.SetCookie("refreshToken", refreshToken, refreshTokenTTLInt, "/auth/refresh-token", cookieDomain, true, true) // chỉ dành cho /auth/refresh-token
		c.SetCookie("refreshToken", refreshToken, refreshTokenTTLInt, "/auth/logout", cookieDomain, true, true)            // chỉ dành cho /auth/logout
		c.SetCookie("refreshToken", refreshToken, refreshTokenTTLInt, "/api/v1/payment/confirm", cookieDomain, true, true) // dành cho /api/v1/payment/confirm
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

// Đăng kí tài khoản
func Register() func(*gin.Context) {
	return func(c *gin.Context) {
		var user models.User // Khoi tao 1 user

		// Kiểm tra nhận được file JSON
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input",
			})
			return
		}

		//fmt.Printf("User received: %+v\n", user)

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

		// Sử dụng DB đã kết nối từ trước (không cần gọi lại ConnectDatabase)
		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Kiểm tra email hoặc username đã tồn tại, chỉ kiểm tra phone nếu có giá trị
		var existingUser models.User
		filter := bson.M{"$or": []bson.M{
			{"username": user.Username},
			{"email": user.Email},
		}}

		// Thêm kiểm tra phone nếu nó được khai báo
		if user.Profile.PhoneNumber != "" {
			filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"profile.phone_number": user.Profile.PhoneNumber})
		}

		err := collection.FindOne(ctx, filter).Decode(&existingUser)
		if err == nil {
			if user.Profile.PhoneNumber != ""{
				c.JSON(http.StatusConflict, gin.H{
					"error": "Email or username or phone already exists.",
				})
				return
			}
			// Nếu không có lỗi, tức là người dùng đã tồn tại
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email or username already exists.",
			})
			return
		}

		// Hash mật khẩu trước khi lưu vào database
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password",
			})
			return
		}
		user.Password = string(hashedPassword)

		// Khoi tao 1 user
		user = newUser(user)

		_, err = collection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
			return
		}

		// Response: 201 Created
		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
			//"user_id": result.InsertedID,
		})
	}
}

// Đăng nhập username/email + password
func Login() func(*gin.Context) {
	return func(c *gin.Context) {
		var loginRequest struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// // Kết nối đến MongoDB
		// if err := config.ConnectDatabase(); err != nil {
		//     c.JSON(http.StatusInternalServerError, gin.H{
		//         "error": "Failed to connect to database",
		//     })
		//     return
		// }

		collection := config.DB.Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Tìm kiếm user bằng email hoặc username
		var user models.User
		filter := bson.M{
			"$or": []bson.M{
				{"email": loginRequest.Username},
				{"username": loginRequest.Username},
			},
		}
		err := collection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Username or password is incorrect",
				}) //Not found with this username
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to find user",
				})
			}
			return
		}

		// Kiểm tra trạng thái is_active
		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Your account has been banned. Please contact support for assistance."})
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
		err = setAuthCookies(c, accessToken, refreshToken, false, true) // set cả 2 cookie
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Trả về đăng nhập thành công và token
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":    accessToken,
			//"refreshToken":    refreshToken,
		})
	}
}

// Đăng xuất
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

// Quên mật khẩu
func ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Email string `json:"email" binding:"required,email"`
		}

		// Kiểm tra dữ liệu đầu vào
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
			})
			return
		}

		// Kết nối đến database và kiểm tra xem người dùng có tồn tại không
		collection := config.DB.Collection("User")
		var user models.User
		err := collection.FindOne(context.TODO(), bson.M{"email": request.Email}).Decode(&user)

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found with this email",
			})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error",
			})
			return
		}

		// // Tạo JWT token
		// token, err := middlewares.GenerateJWT(user.ID.Hex(), user.Role)
		// if err != nil {
		//     c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		//     return
		// }

		// Tạo token ngẫu nhiên và hash nó để lưu vào DB
		rawToken, err := utils.GenerateRandomString(32) // Token ngẫu nhiên
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate reset token",
			})
			return
		}

		hashedToken := utils.HashString(rawToken)    // Hash token để lưu vào DB
		expiresAt := time.Now().Add(3 * time.Minute) // Token có hạn 15 phút

		// Lưu token vào cơ sở dữ liệu
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.M{"_id": user.ID},
			bson.M{
				"$set": bson.M{
					"reset_password_token":   hashedToken,
					"reset_password_expires": expiresAt,
				},
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save reset token",
			})
			return
		}

		// Gửi email
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			log.Fatal("Base URL for reset password is missing")
		}
		resetLink := fmt.Sprintf("%s/reset-password?token=%s", baseURL, rawToken)

		emailTemplatePath := "services/admin_service/templates/password_reset_email.html"
		htmlBody, err := os.ReadFile(emailTemplatePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read email template",
			})
			return
		}

		// Use Go's templating engine to replace the placeholders in the HTML
		t, err := template.New("reset-email").Parse(string(htmlBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to parse email template",
			})
			return
		}

		var bodyBuffer bytes.Buffer
		err = t.Execute(&bodyBuffer, map[string]interface{}{
			"Name":      user.Username,
			"ResetLink": resetLink,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to execute email template",
			})
			return
		}

		//fmt.Println("User Name:", user.Name)
		//fmt.Println("Reset Link:", resetLink)

		htmlBodyString := bodyBuffer.String()

		// Gọi hàm gửi email
		if err := utils.SendEmail(request.Email, "Password Reset Request", htmlBodyString); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to send email",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Password reset link sent to your email",
		})
	}
}

// Đặt lại mật khẩu
func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Token       string `json:"token" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}

		// Kiểm tra dữ liệu đầu vào
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
			})
			return
		}

		// Kiểm tra mật khẩu
		if !utils.IsValidPassword(request.NewPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password must contain at least 8 characters, including letters, numbers, and special characters.",
			})
			return
		}

		// Hash token từ request để so khớp với token trong cơ sở dữ liệu
		hashedToken := utils.HashString(request.Token)

		// Kết nối đến database và tìm người dùng dựa trên token
		userCollection := config.DB.Collection("User")
		var user models.User
		err := userCollection.FindOne(context.TODO(), bson.M{
			"reset_password_token":   hashedToken,
			"reset_password_expires": bson.M{"$gt": time.Now()}, // Chỉ chấp nhận token còn hiệu lực
		}).Decode(&user)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired reset token",
			})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error",
			})
			return
		}

		// Hash mật khẩu mới
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password",
			})
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
		_, err = userCollection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update password",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Password changed successfully",
		})
	}
}

// RefreshToken
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
			"token": accessToken,
		})
	}
}
