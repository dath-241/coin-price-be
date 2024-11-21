package controllers

import (
    "fmt"
	"context"
    "net/http"
    "log"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/config"
    "github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
    "github.com/dath-241/coin-price-be-go/services/admin_service/utils"
    "github.com/gin-gonic/gin"
	"github.com/google/uuid"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// utlis 
func generateUniqueUsername() string {
	return "user-" + uuid.New().String()
}

// Hàm lấy thông tin người dùng từ DB dựa trên email
func getUserByEmail(email string) (*models.User, error) {
    // Lấy collection "users" từ DB
    collection := config.DB.Collection("User")

    var user models.User
    filter := bson.M{"email": email}

    // Truy vấn tìm kiếm người dùng theo email
    err := collection.FindOne(context.TODO(), filter).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            // Nếu không tìm thấy người dùng, trả về lỗi không có tài liệu
            return nil, fmt.Errorf("user not found with email: %s", email)
        }
        log.Println("Error retrieving user:", err)
        return nil, err
    }

    return &user, nil
}

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
    user, err := getUserByEmail(email)
    if err != nil {
        if err.Error() == fmt.Sprintf("user not found with email: %s", email) {
            // Nếu không tìm thấy user, tạo user mới trong DB
            name := userInfo["name"].(string)
			avatarURL := userInfo["picture"].(string)

			var user models.User
			user.Profile.FullName = name
			user.Email = email
			user.Username = generateUniqueUsername()
			user.Profile.AvatarURL = avatarURL

			user = newUser(user)

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

    // Gọi hàm set cookie để thiết lập cookies cho người dùng
    err = setAuthCookies(c, accessToken, refreshToken, false, true) // set cả 2 cookie
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    // Trả về JWT cho frontend
    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "token":   accessToken,
    })
}