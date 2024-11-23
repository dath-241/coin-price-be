package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"strconv"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"strings"
)

func TestAuthMiddleware(t *testing.T) {
	// Mock ACCESS_SECRET
	os.Setenv("ACCESS_SECRET", "mocksecret")

	// Tạo gin router để kiểm tra middleware
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock route có middleware
	router.GET("/protected", AuthMiddleware("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	// Tạo helper để tạo JWT token
	createToken := func(role string, userID string) string {
		claims := jwt.MapClaims{
			"role":    role,
			"user_id": userID,
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
		return tokenString
	}

	// Test case 1: Không có Authorization header
	t.Run("No Authorization Header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Authorization header required")
	})

	// Test case 2: Token không hợp lệ
	t.Run("Invalid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "invalidtoken")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid token")
	})

	// Test case 3: Token hợp lệ nhưng không đủ quyền
	t.Run("Valid Token but Insufficient Role", func(t *testing.T) {
		token := createToken("user", "123")
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusForbidden, resp.Code)
		assert.Contains(t, resp.Body.String(), "Access forbidden: insufficient role")
	})

	// Test case 4: Token hợp lệ và đủ quyền
	t.Run("Valid Token and Sufficient Role", func(t *testing.T) {
		token := createToken("admin", "123")
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Access granted")
	})
}



func TestVerifyJWT(t *testing.T) {
	// Mock environment variables
	os.Setenv("ACCESS_SECRET", "mockaccesssecret")
	os.Setenv("REFRESH_SECRET", "mockrefreshsecret")

	// Helper function to generate token
	createToken := func(secret string, isExpired bool) string {
		expirationTime := time.Now().Add(time.Hour)
		if isExpired {
			expirationTime = time.Now().Add(-time.Hour) // Token đã hết hạn
		}

		claims := &models.CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))
		return tokenString
	}

	// Test case 1: Token hợp lệ (Access Token)
	t.Run("Valid Access Token", func(t *testing.T) {
		token := createToken(os.Getenv("ACCESS_SECRET"), false)
		claims, err := VerifyJWT(token, true)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
	})

	// Test case 2: Token không hợp lệ (sai cấu trúc)
	t.Run("Invalid Token", func(t *testing.T) {
		token := "invalidTokenString"
		claims, err := VerifyJWT(token, true)

		assert.Nil(t, claims)
		assert.EqualError(t, err, "invalid token")
	})

	// Test case 3: Token hết hạn
	t.Run("Expired Token", func(t *testing.T) {
    	token := createToken(os.Getenv("ACCESS_SECRET"), true)
    	claims, err := VerifyJWT(token, true)

	    assert.Nil(t, claims)
    	assert.EqualError(t, err, "invalid token") // Thay đổi kỳ vọng thành "invalid token"
	})


	// Test case 4: Token với signing method không hợp lệ
	t.Run("Invalid Signing Method", func(t *testing.T) {
		claims := &models.CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims) // Sai signing method
		tokenString, _ := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

		result, err := VerifyJWT(tokenString, true)

		assert.Nil(t, result)
		assert.EqualError(t, err, "invalid token")
	})

	// Test case 5: Valid Refresh Token
	t.Run("Valid Refresh Token", func(t *testing.T) {
		token := createToken(os.Getenv("REFRESH_SECRET"), false)
		claims, err := VerifyJWT(token, false)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
	})
}



func TestGenerateAccessToken(t *testing.T) {
	// Thiết lập giá trị mặc định cho biến môi trường
	originalAccessSecret := os.Getenv("ACCESS_SECRET")
	originalAccessTTL := os.Getenv("ACCESS_TOKEN_TTL")
	defer func() {
		os.Setenv("ACCESS_SECRET", originalAccessSecret)
		os.Setenv("ACCESS_TOKEN_TTL", originalAccessTTL)
	}()

	os.Setenv("ACCESS_SECRET", "testsecret") // Set secret key cho test

	t.Run("Success - Generate token", func(t *testing.T) {
		os.Setenv("ACCESS_TOKEN_TTL", "3600") // 1 giờ
		token, err := GenerateAccessToken("12345", "admin")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("testsecret"), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "12345", claims["user_id"])
		assert.Equal(t, "admin", claims["role"])

		exp := int64(claims["exp"].(float64))
		assert.WithinDuration(t, time.Unix(exp, 0), time.Now().Add(3600*time.Second), time.Minute)
	})

	t.Run("Error - ACCESS_TOKEN_TTL not set", func(t *testing.T) {
		os.Unsetenv("ACCESS_TOKEN_TTL") // Xóa biến môi trường
		token, err := GenerateAccessToken("12345", "admin")
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "environment variable ACCESS_TOKEN_TTL is not set")
	})

	t.Run("Error - ACCESS_TOKEN_TTL invalid format", func(t *testing.T) {
		os.Setenv("ACCESS_TOKEN_TTL", "invalid")
		token, err := GenerateAccessToken("12345", "admin")
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "invalid ACCESS_TOKEN_TTL format")
	})

	t.Run("Error - Expired TTL", func(t *testing.T) {
		os.Setenv("ACCESS_TOKEN_TTL", strconv.Itoa(-10)) // TTL âm để token đã hết hạn
		token, err := GenerateAccessToken("12345", "admin")
		assert.NoError(t, err)  // Token vẫn được tạo ra thành công
		assert.NotEmpty(t, token)
	
		// Parse token và kiểm tra lỗi expired
		_, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})
	
		assert.Error(t, err) // Phải có lỗi
		assert.Contains(t, strings.ToLower(err.Error()), "expired") // Kiểm tra chuỗi lỗi chứa "expired"
	})
	
	
}


func TestGenerateRefreshToken(t *testing.T) {
	// Thiết lập các biến môi trường cần thiết
	os.Setenv("REFRESH_SECRET", "testrefreshsecret")

	t.Run("Success - Generate token", func(t *testing.T) {
		os.Setenv("REFRESH_TOKEN_TTL", "3600") // TTL 1 giờ
		token, err := GenerateRefreshToken("12345", "user")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Parse token để kiểm tra nội dung
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("testrefreshsecret"), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "12345", claims["user_id"])
		assert.Equal(t, "user", claims["role"])
		assert.NotEmpty(t, claims["exp"])
	})

	t.Run("Error - REFRESH_TOKEN_TTL not set", func(t *testing.T) {
		os.Unsetenv("REFRESH_TOKEN_TTL") // Xóa biến môi trường
		token, err := GenerateRefreshToken("12345", "user")
		assert.Error(t, err)
		assert.Equal(t, "", token)
		assert.Contains(t, err.Error(), "REFRESH_TOKEN_TTL is not set")
	})

	t.Run("Error - REFRESH_TOKEN_TTL invalid format", func(t *testing.T) {
		os.Setenv("REFRESH_TOKEN_TTL", "invalid") // TTL không hợp lệ
		token, err := GenerateRefreshToken("12345", "user")
		assert.Error(t, err)
		assert.Equal(t, "", token)
		assert.Contains(t, err.Error(), "invalid REFRESH_TOKEN_TTL format")
	})

	t.Run("Error - Expired TTL", func(t *testing.T) {
		os.Setenv("REFRESH_TOKEN_TTL", strconv.Itoa(-10)) // TTL âm để token hết hạn
		token, err := GenerateRefreshToken("12345", "user")
		assert.NoError(t, err) // Token vẫn được tạo thành công
		assert.NotEmpty(t, token)

		// Parse token để kiểm tra lỗi expired
		_, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("testrefreshsecret"), nil
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expired")
	})
}