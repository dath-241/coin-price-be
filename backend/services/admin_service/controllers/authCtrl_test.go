package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func TestSetAuthCookies(t *testing.T) {
	
	gin.SetMode(gin.TestMode)

	// Tạo trường hợp test
	tests := []struct {
		name             string
		setup            func()
		accessToken      string
		refreshToken     string
		setAccessToken   bool
		setRefreshToken  bool
		expectedError    string
		expectedCookies  []http.Cookie
	}{
		{
			name:            "Test case 1: Set both access and refresh tokens",
			accessToken:     "access_token_value",
			refreshToken:    "refresh_token_value",
			setAccessToken:  true,
			setRefreshToken: true,
			expectedError:   "",
			expectedCookies: []http.Cookie{
				{Name: "accessToken", Value: "access_token_value", Path: "/api/v1", Domain: "example.com", HttpOnly: true, Secure: true},
				{Name: "refreshToken", Value: "refresh_token_value", Path: "/api/v1/auth/refresh-token", Domain: "example.com", HttpOnly: true, Secure: true},
			},
			setup: func() {
				os.Setenv("COOKIE_DOMAIN", "example.com")
				os.Setenv("ACCESS_TOKEN_TTL", "3600")  // 1 giờ
				os.Setenv("REFRESH_TOKEN_TTL", "7200") // 2 giờ
			},
		},
		{
			name:            "Testcase 2: Missing environment variables",
			accessToken:     "access_token_value",
			refreshToken:    "refresh_token_value",
			setAccessToken:  true,
			setRefreshToken: true,
			expectedError:   "environment variables are not set",
			setup: func() {
				os.Unsetenv("COOKIE_DOMAIN")
				os.Unsetenv("ACCESS_TOKEN_TTL")
				os.Unsetenv("REFRESH_TOKEN_TTL")
			},
		},
		{
			name:            "Test case 3: Set access token",
			accessToken:     "access_token_value",
			refreshToken:    "refresh_token_value",
			setAccessToken:  true,
			setRefreshToken: false,
			expectedError:   "",
			expectedCookies: []http.Cookie{
				{Name: "accessToken", Value: "access_token_value", Path: "/api/v1", Domain: "example.com", HttpOnly: true, Secure: true},
			},
			setup: func() {
				os.Setenv("COOKIE_DOMAIN", "example.com")
				os.Setenv("ACCESS_TOKEN_TTL", "3600")  // 1 giờ
				os.Setenv("REFRESH_TOKEN_TTL", "7200") // 2 giờ
			},
		},
		{
			name:            "Test case 4: Set refresh token",
			accessToken:     "access_token_value",
			refreshToken:    "refresh_token_value",
			setAccessToken:  false,
			setRefreshToken: true,
			expectedError:   "",
			expectedCookies: []http.Cookie{
				{Name: "refreshToken", Value: "refresh_token_value", Path: "/api/v1/auth/refresh-token", Domain: "example.com", HttpOnly: true, Secure: true},
			},
			setup: func() {
				os.Setenv("COOKIE_DOMAIN", "example.com")
				os.Setenv("ACCESS_TOKEN_TTL", "3600")  // 1 giờ
				os.Setenv("REFRESH_TOKEN_TTL", "7200") // 2 giờ
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Gọi setup nếu được định nghĩa
			if tt.setup != nil {
				tt.setup()
			}

			// Tạo ngữ cảnh Gin và trình ghi HTTP
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Gọi hàm setAuthCookies
			err := setAuthCookies(c, tt.accessToken, tt.refreshToken, tt.setAccessToken, tt.setRefreshToken)

			// Kiểm tra lỗi
			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
				//fmt.Println("Error:", err) // Debug lỗi nếu không khớp

				// Đảm bảo không có cookie nào được thiết lập
				assert.Empty(t, w.Result().Cookies(), "Cookies should not be set when error")
			} else {
				assert.NoError(t, err)

				// Kiểm tra cookie trong phản hồi
				for _, expectedCookie := range tt.expectedCookies {
					found := false
					for _, cookie := range w.Result().Cookies() {
						if cookie.Name == expectedCookie.Name &&
							cookie.Value == expectedCookie.Value &&
							cookie.Path == expectedCookie.Path &&
							cookie.Domain == expectedCookie.Domain &&
							cookie.HttpOnly == expectedCookie.HttpOnly &&
							cookie.Secure == expectedCookie.Secure {
							found = true
						}
					}
					assert.True(t, found, "Cookie %s not found", expectedCookie.Name)
				}
			}
		})
	}
}

func TestResetAuthCookies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		setup              func() // Hàm setup cho các test riêng biệt
		expectedError      string
		expectedCookies    []http.Cookie // Kiểm tra xem có cookies được thiết lập hay không
	}{
		{
			name: "Test case 1: Reset cookies successfully",
			setup: func() {
				// Set up môi trường và cookies cho test này
				os.Setenv("COOKIE_DOMAIN", "example.com")
			},
			expectedError:   "",
			expectedCookies: []http.Cookie{
				{Name: "accessToken", Value: "", Path: "/", Domain: "example.com", MaxAge: 0, HttpOnly: true, Secure: true},
                {Name: "refreshToken", Value: "", Path: "/", Domain: "example.com", MaxAge: 0, HttpOnly: true, Secure: true},
			}, 
		},
		{
			name: "Test case 2: Missing environment variables",
			setup: func() {
				// Không thiết lập biến môi trường COOKIE_DOMAIN để kiểm tra lỗi
				os.Unsetenv("COOKIE_DOMAIN")
			},
			expectedError: "environment variables are not set",
			expectedCookies: nil, // Cookies không được thiết lập vì có lỗi môi trường
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Gọi hàm setup nếu có
			if tt.setup != nil {
				tt.setup()
			}

			// Tạo ngữ cảnh Gin
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Gọi hàm resetAuthCookies
			err := resetAuthCookies(c)

			// Kiểm tra lỗi nếu có
			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Fatalf("expected error %v, got %v", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			// Kiểm tra cookies trong phản hồi
			for _, expectedCookie := range tt.expectedCookies {
				found := false
				for _, cookie := range w.Result().Cookies() {
					if cookie.Name == expectedCookie.Name &&
						cookie.Value == expectedCookie.Value && // Kiểm tra giá trị cookie là rỗng
						cookie.Path == expectedCookie.Path &&
						cookie.Domain == expectedCookie.Domain &&
						cookie.HttpOnly == expectedCookie.HttpOnly &&
						cookie.Secure == expectedCookie.Secure {
						found = true
					}
				}
				if !found {
					t.Errorf("Cookie %s not found or not cleared properly", expectedCookie.Name)
				}
			}
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

	// // Test case 1: Đăng nhập thành công
	// t.Run("Login successful", func(t *testing.T) {
	// 	loginData := map[string]string{
	// 		"username": "testuser",
	// 		"password": "password123",
	// 	}
	// 	body, _ := json.Marshal(loginData)

	// 	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	// 	req.Header.Set("Content-Type", "application/json")
	// 	w := httptest.NewRecorder()

	// 	r.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusOK, w.Code)

	// 	var response map[string]interface{}
	// 	err := json.Unmarshal(w.Body.Bytes(), &response)
	// 	assert.NoError(t, err)

	// 	assert.Equal(t, "Login successful", response["message"])
	// 	assert.NotEmpty(t, response["token"])
	// })

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
		accessToken        string
		refreshToken       string
		expectedStatusCode int
		expectedBody       string
	}{
		// {
		// 	name:               "Logout successful with valid tokens",
		// 	accessToken:        "validAccessToken",
		// 	refreshToken:       "validRefreshToken",
		// 	expectedStatusCode: http.StatusOK,
		// 	expectedBody:       `{"message":"Logout successful"}`,
		// },
		{
			name:               "No access token provided",
			accessToken:        "",
			refreshToken:       "validRefreshToken",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"No token provided"}`,
		},
		{
			name:               "No refresh token provided",
			accessToken:        "validAccessToken",
			refreshToken:       "",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Refresh Token not provided"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("MONGO_URI", "mongodb://localhost:27017")
			os.Setenv("MONGO_DB_NAME", "test_db")
			os.Setenv("ACCESS_TOKEN_TTL", "3600")
			os.Setenv("REFESH_TOKEN_TTL", "3600")

			// Setup Gin
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/logout", Logout())


			// Create a mock request
			req := httptest.NewRequest(http.MethodPost, "/logout", nil)
			if tt.accessToken != "" {
				req.Header.Set("Authorization", tt.accessToken)
			}
			if tt.refreshToken != "" {
				req.AddCookie(&http.Cookie{
					Name:  "refreshToken",
					Value: tt.refreshToken,
				})
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
			defer os.Unsetenv("ACCESS_TOKEN_TTL")
			defer os.Unsetenv("REFRESH_TOKEN_TTL")
		})
	}
}

