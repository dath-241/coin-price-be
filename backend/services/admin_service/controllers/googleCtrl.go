package controllers

import (
    "fmt"
	"context"
    "net/http"
    "strconv"
    "os"

	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/config"
    "github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
    "github.com/dath-241/coin-price-be-go/services/admin_service/utils"
    "github.com/gin-gonic/gin"
)

// GoogleLogin xử lý đăng nhập bằng Google ID Token
func GoogleLogin(c *gin.Context) {
    idToken := c.PostForm("id_token") // Nhận Google ID Token từ frontend

    // Xác minh Google ID Token
    userInfo, err := utils.VerifyGoogleIDToken(idToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Invalid Google ID token",
        })
        return
    }

    // userID := userInfo["sub"].(string) // Lấy user ID từ token
    // Tạo email từ thông tin Google
    email := userInfo["email"].(string)
    
    // Lấy thông tin user từ DB
    user, err := utils.GetUserByEmail(email)
    if err != nil {
        if err.Error() == fmt.Sprintf("user not found with email: %s", email) {
            // Nếu không tìm thấy user, tạo user mới trong DB
            name := userInfo["name"].(string)
            user = &models.User{
                //ID: ID,
                Name: name,
                Email:  email,
                Role:   "VIP-0", // Gán role mặc định cho user mới, có thể là "user" hoặc một role khác
                // Thêm các trường khác nếu cần
            }
            collection := config.DB.Collection("User")
            _, insertErr := collection.InsertOne(context.TODO(), user)
            if insertErr != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Failed to create new user",
                })
                return
            }
        } else {
            // Lỗi khác trong quá trình truy xuất DB
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Error retrieving user from database",
            })
            return
        }
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

    // Load biến môi trường cho tên miền cookie và thời gian sống
    cookieDomain := os.Getenv("COOKIE_DOMAIN")
    accessTokenTTL := os.Getenv("ACCESS_TOKEN_TTL")
    refreshTokenTTL := os.Getenv("REFRESH_TOKEN_TTL")

    if cookieDomain == "" || accessTokenTTL == "" || refreshTokenTTL == "" {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Environment variables are not set",
        })
        return
    }

    accessTokenTTLInt, err := strconv.Atoi(accessTokenTTL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Invalid ACCESS_TOKEN_TTL format",
        })
        return
    }

    refreshTokenTTLInt, err := strconv.Atoi(refreshTokenTTL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Invalid REFRESH_TOKEN_TTL format",
        })
        return
    }

    // Gửi token dưới dạng cookie
    // Cookie cho xác thực
    c.SetCookie("accessToken", accessToken, accessTokenTTLInt, "/api/v1", cookieDomain, true, true)  // chỉ dành cho /api/v1
    c.SetCookie("refreshToken", refreshToken, refreshTokenTTLInt, "/auth/refresh-token", cookieDomain, true, true) // chỉ dành cho /auth/refresh-token

    // Cookie cho các hành động logout hoặc các route riêng biệt
    c.SetCookie("accessToken", accessToken, accessTokenTTLInt, "/auth/logout", cookieDomain, true, true)  // chỉ dành cho /auth/logout
    c.SetCookie("refreshToken", refreshToken, refreshTokenTTLInt, "/auth/logout", cookieDomain, true, true) // chỉ dành cho /auth/logout
    c.SetCookie("refreshToken", refreshToken, refreshTokenTTLInt, "/api/v1/payment/confirm", cookieDomain, true, true) // dành cho /api/v1/payment/confirm

    // Trả về JWT cho frontend
    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        //"token":   token,
    })
}

