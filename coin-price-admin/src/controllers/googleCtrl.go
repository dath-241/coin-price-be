package controllers

import (
    "fmt"
	"context"
    "net/http"

	"coin-price-admin/src/models"
	"coin-price-admin/src/config"
    "coin-price-admin/src/middlewares"
    "coin-price-admin/src/utils"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// // Xác minh Google ID Token bằng OAuth2 Client
// func verifyGoogleIDToken(idToken string) (map[string]interface{}, error) {
//     // Sử dụng Google OAuth2 client để xác minh ID Token
//     ctx := context.Background()
//     client, err := oauth2.NewService(ctx, option.WithCredentialsFile("path_to_credentials.json"))
//     if err != nil {
//         return nil, fmt.Errorf("could not create OAuth2 service: %v", err)
//     }

//     // Kiểm tra token qua Google API
//     tokenInfo, err := client.Tokeninfo().IdToken(idToken).Do()
//     if err != nil {
//         return nil, fmt.Errorf("invalid token: %v", err)
//     }

//     return map[string]interface{}{
//         "sub":   tokenInfo.UserId,
//         "email": tokenInfo.Email,
//     }, nil
// }

// Hàm tạo JWT


// GoogleLogin xử lý đăng nhập bằng Google ID Token
func GoogleLogin(c *gin.Context) {
    idToken := c.PostForm("id_token") // Nhận Google ID Token từ frontend

    // Xác minh Google ID Token
    userInfo, err := utils.VerifyGoogleIDToken(idToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google ID token"})
        return
    }

    
    // userID := userInfo["sub"].(string) // Lấy user ID từ token
    // Tạo email từ thông tin Google
    email := userInfo["email"].(string)
    ID:= primitive.NewObjectID()
    // Lấy thông tin user từ DB
    user, err := utils.GetUserByEmail(email)
    if err != nil {
        if err.Error() == fmt.Sprintf("user not found with email: %s", email) {
            // Nếu không tìm thấy user, tạo user mới trong DB
            user = &models.User{
                ID: ID,
                Email:  email,
                Role:   "VIP-0", // Gán role mặc định cho user mới, có thể là "user" hoặc một role khác
                // Thêm các trường khác nếu cần
            }
            collection := config.DB.Collection("User")
            _, insertErr := collection.InsertOne(context.TODO(), user)
            if insertErr != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new user"})
                return
            }
        } else {
            // Lỗi khác trong quá trình truy xuất DB
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user from database"})
            return
        }
    }

    // Lấy role từ user (giả sử role đã được lưu trong DB)
    role := user.Role

    // Tạo JWT cho người dùng
    token, err := middlewares.GenerateJWT(user.ID.Hex(), role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
        return
    }

    // // Lưu JWT vào cookie (hoặc trả lại token cho frontend)
    // c.SetCookie("token", token, 3600, "/", "localhost", false, true)

    // Trả về JWT cho frontend
    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "token":   token,
    })
}

