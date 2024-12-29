package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	//"strings"
	"testing"

	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	// Bảng các test cases
	tests := []struct {
		name           		string
		user           		models.User
		expectedUsername	string
		expectedEmail		string
		expectedRole   		string
		expectedActive 		bool
		expectedFullName	string
		//expectedPhoneNumber	string
		//expectedDateOfBirth	string
		expectedAvatar 		string
		//expectedBio			string
	}{
		{
			name: "Test case 1: New User with default values",
			user: models.User{
				Username: "testuser1",
				Email: "quangtac2004@gmail.com",
				//Profile:  models.Profile{},
			},
			expectedUsername: "testuser1",
			expectedEmail: "quangtac2004@gmail.com",
			expectedRole:   "VIP-0",  // Role mặc định
			expectedAvatar: "https://drive.google.com/file/d/15Ef4yebpGhT8pwgnt__utSESZtJdmA4a/view?usp=sharing", // Avatar mặc định
			expectedActive: true,     // Trạng thái người dùng mặc định
			expectedFullName: "testuser1",
		},
		{
			name: "Test case 2: User with custom profile",
			user: models.User{
				Username: "testuser2",
				Email: "quangtac2004@gmail.com",
				Profile: models.Profile{
					FullName: "Test User 2",
					AvatarURL: "https://example.com/custom-avatar.jpg",
				},
			},
			expectedUsername: "testuser2",
			expectedEmail: "quangtac2004@gmail.com",
			expectedRole:   "VIP-0",
			expectedAvatar: "https://example.com/custom-avatar.jpg", // Avatar được chỉ định
			expectedActive: true,
			expectedFullName: "Test User 2",
		},
		// Thêm nhiều test case khác nếu cần
	}

	// Duyệt qua tất cả các test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newUser(tt.user)

			// Kiểm tra các giá trị kỳ vọng
			assert.Equal(t, tt.expectedUsername, result.Username)
			assert.Equal(t, tt.expectedEmail, result.Email)
			assert.Equal(t, tt.expectedRole, result.Role)
			assert.Equal(t, tt.expectedActive, result.IsActive)
			assert.Equal(t, tt.expectedAvatar, result.Profile.AvatarURL)
			
			assert.NotNil(t, result.CreatedAt)
			assert.NotNil(t, result.UpdatedAt)
		})
	}
}


func TestRegister(t *testing.T) {
	mockRepo := &repository.MockUserRepository{
		Users: make(map[string]interface{}),
	}

	// Mock router
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", Register(mockRepo))

	t.Run("Valid Registration", func(t *testing.T) {
		newUser := models.User{
			Username: "valid-user",
			Email:    "valid@example.com",
			Password: "Valid@1234",
		}

		reqBody, _ := json.Marshal(newUser)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Perform request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"message": "User registered successfully"}`, w.Body.String())
	})

	t.Run("Invalid Username", func(t *testing.T) {
		invalidUser := models.User{
			Username: "invalid username",
			Email:    "valid@example.com",
			Password: "Valid@1234",
		}

		reqBody, _ := json.Marshal(invalidUser)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Perform request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Username only alphanumeric characters and hyphens are allowed."}`, w.Body.String())
	})

	t.Run("Weak Password", func(t *testing.T) {
		invalidUser := models.User{
			Username: "valid-user",
			Email:    "valid@example.com",
			Password: "1234",
		}

		reqBody, _ := json.Marshal(invalidUser)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Perform request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Password must contain at least 8 characters, including letters, numbers, and special characters."}`, w.Body.String())
	})
}

func TestLogin(t *testing.T) {
	// Mock dữ liệu người dùng
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUsers := map[string]interface{}{
		"12345": models.User{
			ID:       primitive.NewObjectID(),
			Username: "testuser",
			Email:    "testuser@example.com",
			Password: string(passwordHash),
			IsActive: true,
			Role:     "user",
		},
	}

	// Mock repository
	mockRepo := &repository.MockUserRepository{
		Users: mockUsers,
	}

	// Khởi tạo Gin router và handler
	r := gin.Default()
	r.POST("/login", Login(mockRepo))

	// Test case 2: Sai mật khẩu
	t.Run("Incorrect password", func(t *testing.T) {
		loginData := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(loginData)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Username or password is incorrect", response["error"])
	})
}

func TestLogout(t *testing.T) {

	tests := []struct {
		name               string
		tokenString        string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "No token provided",
			tokenString:        "",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"No token provided"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("MONGO_URI", "mongodb://localhost:27017")
			os.Setenv("MONGO_DB_NAME", "test_db")
			os.Setenv("JWT_TOKEN_TTL", "3600")

			// Setup Gin
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/logout", Logout())


			// Create a mock request
			req := httptest.NewRequest(http.MethodPost, "/logout", nil)
			if tt.tokenString != "" {
				req.Header.Set("Authorization", tt.tokenString)
			}

			// Create a mock response recorder
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the status code and body
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			
			defer os.Unsetenv("MONGO_URI")
			defer os.Unsetenv("MONGO_DB_NAME")
			defer os.Unsetenv("JWT_TOKEN_TTL")
		})
	}
}

