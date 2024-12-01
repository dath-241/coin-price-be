package middlewares 

import (
	"os"
    "fmt"
	"time"
    "errors"
    "strconv"
    "net/http"
    "sync"
    "github.com/dath-241/coin-price-be-go/services/admin_service/models"
	
    "github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

)

var (
	BlacklistedTokens      = make(map[string]time.Time) // Token và thời gian hết hạn
	BlacklistedTokensMutex sync.Mutex                  // Mutex để đồng bộ
)

// Hàm kiểm tra phân quyền token người dùng
func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")

        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Kiểm tra xem token có trong danh sách từ chối và hết hạn chưa
        if expTime, found := BlacklistedTokens[tokenString]; found {
            if time.Now().After(expTime) {
                delete(BlacklistedTokens, tokenString) // Xóa token đã hết hạn khỏi danh sách từ chối
            } else {
                c.JSON(http.StatusUnauthorized, gin.H{
                    "error": "Token has been revoked",
                })
                c.Abort()
                return
            }
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil // JWT_SECRET
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token",
            })
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token claims",
            })
            c.Abort()
            return
        }

        userRole := claims["role"].(string)
        
        // Kiểm tra quyền
        hasAccess := false
        for _, role := range allowedRoles {
            if userRole == role {
                hasAccess = true
                break
            }
        }

        if !hasAccess {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Access forbidden: insufficient role",
            })
            c.Abort()
            return
        }

        // Cho phép tiếp tục nếu xác thực thành công
        c.Set("user_id", claims["user_id"])
        c.Next()
    }
}

// VerifyJWT sẽ xác thực JWT và trả về các claims nếu token hợp lệ
func VerifyJWT(tokenString string) (*models.CustomClaims, error) {
    // Lấy secret key từ environment
    var jwtKey = []byte(os.Getenv("JWT_SECRET"))

    // Parse và kiểm tra token
    token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtKey, nil
    })

    // Kiểm tra lỗi khi parse token
    if err != nil {
        return nil, errors.New("invalid token")
    }

    // Kiểm tra tính hợp lệ của các claims
    if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
        if claims.ExpiresAt.Time.Before(time.Now()) {
            return nil, errors.New("token has expired")
        }
        return claims, nil
    }

    return nil, errors.New("invalid token claims")
}



func GenerateToken(userID, role string) (string, error) {
    tokenSecret := []byte(os.Getenv("JWT_SECRET")) // Lấy khóa bí mật từ biến môi trường
    
    // Load biến môi trường cho thời gian sống của access token
    tokenTTL := os.Getenv("JWT_TOKEN_TTL")
    if tokenTTL == "" {
        return "", fmt.Errorf("environment variable JWT_TOKEN_TTL is not set")
    }

    tokenTTLInt, err := strconv.Atoi(tokenTTL)
    if err != nil {
        return "", fmt.Errorf("invalid JWT_TOKEN_TTL format: %v", err)
    }

    claims := jwt.MapClaims{
        "user_id":  userID,
        "role":     role,
        "exp":      time.Now().Add(time.Duration(tokenTTLInt) * time.Second).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(tokenSecret)
}