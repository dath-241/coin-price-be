package controllers 

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
    "log"
    "bytes"
    "html/template"
	"backend/services/admin_service/src/config"
	"backend/services/admin_service/src/models"
    "backend/services/admin_service/src/middlewares"
    "backend/services/admin_service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

    
// Đăng kí tài khoản
func Register() func(*gin.Context) {
    return func(c *gin.Context) {
        var user models.User    // Khoi tao 1 user

        // Kiểm tra nhận được file JSON
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }

        // Kiểm tra xem tên, email và mật khẩu có null không
        if user.Name == "" || user.Email == "" || user.Password == "" {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Name, email and password are required.",
            })
            return
        }
        // Kiểm tra độ dài tên
        if !utils.IsValidName(user.Name) {
            c.JSON(http.StatusBadRequest, gin.H{
            "error": "Name length must be between 1 and 50 characters.",
            })
            return
        }
            
        // Kiểm tra xem tên chỉ chứa các ký tự chữ cái
        if !utils.IsAlphabetical(user.Name) {
            c.JSON(http.StatusBadRequest, gin.H{
            "error": "Name must only contain alphabetical characters.",
            })
            return
        }
        // Kiểm tra mật khẩu
        if !utils.IsValidPassword(user.Password) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least 8 characters, including letters, numbers, and special characters."})
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

        // Kiểm tra email đã tồn tại
        var existingUser models.User
        err := collection.FindOne(ctx, bson.M{"$or": []bson.M{
            //{"username": user.Username},
            {"email": user.Email},
        }}).Decode(&existingUser)

        if err == nil {
            // Nếu không có lỗi, tức là người dùng đã tồn tại
            c.JSON(http.StatusConflict, gin.H{
                "error": "Email already exists.",
            })
            return
        }

        user.Role = "VIP-0"
        user.CreatedAt = time.Now()

        // Hash mật khẩu trước khi lưu vào database
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
            return
        }
        user.Password = string(hashedPassword)

        result, err := collection.InsertOne(ctx, user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
            return
        }

        // Response: 201 Created
        c.JSON(http.StatusCreated, gin.H{
            "message":    "User registered successfully",
            "user_id":    result.InsertedID,
        })
    }
}

func Login() func(*gin.Context) {
    return func(c *gin.Context) {
        var loginRequest struct {
            Username    string `json:"username" binding:"required"`
            Password    string `json:"password" binding:"required"`
        }

        if err := c.ShouldBindJSON(&loginRequest); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Kết nối đến MongoDB
        if err := config.ConnectDatabase(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
            return
        }

        collection := config.DB.Collection("User")

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        // Kiểm tra tài khoản có tồn tại 
        var user models.User
        err := collection.FindOne(ctx, bson.M{"email": loginRequest.Username}).Decode(&user)
        if err != nil {
            if err == mongo.ErrNoDocuments {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Username or password is incorrect"}) //Not found with this username
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
            }
            return
        }

        // Kiểm tra mật khẩu
        err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Username or password is incorrect"})
            return
        }

        // Tạo JWT token
        token, err := middlewares.GenerateJWT(user.ID.Hex(), user.Role)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
            return
        }

        // Trả về đăng nhập thành công và token
        c.JSON(http.StatusOK, gin.H{
            "message": "Login successful",
            "token":    token,
        })
    }
}

// Hàm Logout
func Logout() func(*gin.Context) {
    return func(c *gin.Context) {
        // Lấy token từ header
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
            return
        }

        // Parse token và lấy thời gian hết hạn
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        // Lấy claims từ token
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            return
        }

        // Lấy thời gian hết hạn từ claims
        expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

        // Thêm token vào danh sách Blacklisted với thời gian hết hạn
        middlewares.BlacklistedTokens[tokenString] = expirationTime

        // Trả về thông báo thành công
        c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
    }
}


// Đặt lại mật khẩu
func ForgotPassword() gin.HandlerFunc {
    return func(c *gin.Context) {
        var request struct {
            Email string `json:"email" binding:"required,email"`
        }

        // Kiểm tra dữ liệu đầu vào
        if err := c.ShouldBindJSON(&request); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
            return
        }

        // Kết nối đến database và kiểm tra xem người dùng có tồn tại không
        collection := config.DB.Collection("User")
        var user models.User
        err := collection.FindOne(context.TODO(), bson.M{"email": request.Email}).Decode(&user)
        
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found with this email"})
            return
        } else if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            return
        }

        // Tạo JWT token
        token, err := middlewares.GenerateJWT(user.ID.Hex(), user.Role)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
            return
        }

        // Gửi email
        baseURL := os.Getenv("BASE_URL")
        if baseURL == "" {
            log.Fatal("Base URL for reset password is missing")
        }
        resetLink := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)

        emailTemplatePath := "templates/password_reset_email.html"
        htmlBody, err := os.ReadFile(emailTemplatePath)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read email template"})
            return
        }

        // Use Go's templating engine to replace the placeholders in the HTML
        t, err := template.New("reset-email").Parse(string(htmlBody))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse email template"})
            return
        }

        var bodyBuffer bytes.Buffer
        err = t.Execute(&bodyBuffer, map[string]interface{}{
            "Name":      user.Name,
            "ResetLink": resetLink,
        })
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute email template"})
            return
        }

        fmt.Println("User Name:", user.Name)
        fmt.Println("Reset Link:", resetLink)   

        htmlBodyString := bodyBuffer.String()
       
        // Gọi hàm gửi email
        if err := utils.SendEmail(request.Email, "Password Reset Request", htmlBodyString); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Reset password email sent"})
    }
}


func ResetPassword() gin.HandlerFunc {
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

        // Kiểm tra mật khẩu
        if !utils.IsValidPassword(request.NewPassword) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least 8 characters, including letters, numbers, and special characters."})
            return
        }

        // Giải mã và xác thực token
        claims, err := middlewares.VerifyJWT(request.Token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            return
        }


        // Kết nối đến database và kiểm tra xem người dùng có tồn tại không
        userID := claims.UserID // Lấy ID người dùng từ claims

        // Chuyển đổi userID từ chuỗi sang ObjectID
        objID, err := primitive.ObjectIDFromHex(userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
            return
        }

        
        userCollection := config.DB.Collection("User")
        var user models.User
        err = userCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found with this email"})
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

        // Cập nhật mật khẩu mới cho người dùng và xóa OTP
        update := bson.M{"$set": bson.M{"password": hashedPassword}, "$unset": bson.M{"otp": ""}} // Xóa OTP sau khi sử dụng
        _, err = userCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
    }
}


func RefreshToken() func(*gin.Context) {
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

        // Lấy user_id và role từ claims của token cũ
        userID := claims["user_id"].(string)
        userRole := claims["role"].(string)

        // Lấy thời gian hết hạn từ claims
        expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

        // Thêm token cũ vào blacklist
        middlewares.BlacklistedTokens[tokenString] = expirationTime

        // Tạo JWT token mới cho người dùng
        newToken, err := middlewares.GenerateJWT(userID, userRole)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
            return
        }

        // Trả về token mới
        c.JSON(http.StatusOK, gin.H{
            "token":   newToken,
        })
    }
}
