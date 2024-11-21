package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	//"strings"
	"testing"
	"time"

	//"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/utils"
	"github.com/gin-gonic/gin"

	//"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// MockMongoCollection mô phỏng MongoDB Collection
type MockMongoCollection struct {
	mock.Mock
}

func (m *MockMongoCollection) FindOne(ctx context.Context, filter interface{}) *MockSingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*MockSingleResult)
}

// Mock InsertOne method for mocking insert operation
func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// Mock UpdateOne for mocking update operation
func (m *MockMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
    args := m.Called(ctx, filter, update)
    return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

// MockSingleResult mô phỏng một SingleResult để decode
type MockSingleResult struct {
	mock.Mock
}

// func (m *MockSingleResult) Decode(v interface{}) error {
// 	args := m.Called(v)
// 	if user, ok := v.(*models.User); ok { // Đảm bảo v là con trỏ tới models.User
// 		*user = args.Get(0).(models.User) // Gán giá trị vào con trỏ
// 	}
// 	return args.Error(1) // Trả về lỗi nếu có
// }

func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	if user, ok := v.(*models.User); ok {
		// Kiểm tra args.Get(0) có phải là nil hay không
		if res := args.Get(0); res != nil {
			*user = res.(models.User) // Gán giá trị vào con trỏ nếu không phải nil
		} else {
			// Xử lý khi không có dữ liệu trả về
			return fmt.Errorf("no data found for Decode")
		}
	}
	return args.Error(1) // Trả về lỗi nếu có
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	mockCollection := new(MockMongoCollection)

	r.POST("/login", func(c *gin.Context) {
		var loginRequest struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := mockCollection.FindOne(context.TODO(), map[string]interface{}{
			"username": loginRequest.Username,
		})

		var user models.User
		if err := result.Decode(&user); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Username or password is incorrect"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Username or password is incorrect"})
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Your account has been banned. Please contact support for assistance."})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
		})
	})

	// Test case 1: Đăng nhập thành công
	t.Run("Test case 1: success", func(t *testing.T) {
		
		// Mock data cho test case success
		userID := primitive.NewObjectID()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("dath2024@"), bcrypt.DefaultCost)
		validUser := models.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
			Password: string(hashedPassword),
			IsActive: true,
			Role:     "admin",
		}
		mockResult := new(MockSingleResult)
		// Mock logic cho FindOne và Decode
		mockResult.On("Decode", mock.Anything).Return(validUser, nil)
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockResult)

		// Tạo request
		loginRequest := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: "testuser",
			Password: "dath2024@",
		}
		body, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		// Gửi request và kiểm tra kết quả
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Login successful")

		// Kiểm tra expectation của mock
		mockResult.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
	})

	// Test case 2: User không tồn tại
	t.Run("Test case 2: user not found", func(t *testing.T) {
		
		mockCollection = new(MockMongoCollection)
		mockResult := new(MockSingleResult)

		// Mock logic cho user not found
		mockResult.On("Decode", mock.Anything).Return(models.User{}, assert.AnError)
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockResult)

		// Tạo request
		loginRequest := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: "wronguser",
			Password: "password",
		}
		body, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		// Gửi request và kiểm tra kết quả
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Username or password is incorrect")

		// Kiểm tra expectation của mock
		mockResult.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
	})
	
	// Test case 3: Sai mật khẩu
	t.Run("Test case 3: wrong password", func(t *testing.T) {
		mockCollection = new(MockMongoCollection)
		mockResult := new(MockSingleResult)

		// Mock dữ liệu cho người dùng hợp lệ nhưng mật khẩu không chính xác
		userID := primitive.NewObjectID()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
		validUser := models.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
			Password: string(hashedPassword),
			IsActive: true,
			Role:     "admin",
		}

		// Mock logic cho FindOne và Decode
		// Mock logic cho Decode
		mockResult.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
			// Kiểm tra kiểu của args.Get(0)
			userPtr, ok := args.Get(0).(*models.User)
			if !ok {
				t.Error("Expected *models.User but got a different type")
				return
			}
			// Gán giá trị đúng vào đối tượng user
			*userPtr = validUser
		}).Return(nil)

		mockCollection.On("FindOne", mock.Anything, map[string]interface{}{
			"username": "testuser",
		}).Return(mockResult)

		// Tạo request với mật khẩu sai
		loginRequest := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: "testuser",
			Password: "wrong_password", // Mật khẩu không chính xác
		}
		body, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		// Gửi request và kiểm tra phản hồi
		r.ServeHTTP(w, req)

		// Kỳ vọng trả về mã lỗi 401 và thông báo phù hợp
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Username or password is incorrect")

		// Kiểm tra expectation của mock
		mockResult.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
	})
}

func TestRegister(t *testing.T) {
    gin.SetMode(gin.TestMode)

    r := gin.Default()
    mockCollection := new(MockMongoCollection)

    r.POST("/register", func(c *gin.Context) {
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

        // Mock DB logic kiểm tra user tồn tại
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        var existingUser models.User
        filter := bson.M{"$or": []bson.M{
            {"username": user.Username},
            {"email": user.Email},
        }}

        if user.Profile.PhoneNumber != "" {
            filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"profile.phone_number": user.Profile.PhoneNumber})
        }

        // Mock find result
        mockCollection.On("FindOne", ctx, filter).Return(mockCollection)

        err := mockCollection.FindOne(ctx, filter).Decode(&existingUser)
        if err == nil {
            c.JSON(http.StatusConflict, gin.H{
                "error": "Email or username or phone already exists.",
            })
            return
        }

        // Mock hash mật khẩu
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to hash password",
            })
            return
        }
        user.Password = string(hashedPassword)

        // Mock logic tạo user mới trong DB
        _, err = mockCollection.InsertOne(ctx, user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to create user",
            })
            return
        }

        c.JSON(http.StatusCreated, gin.H{
            "message": "User registered successfully",
        })
    })

    // Test case: Đăng ký thành công
	t.Run("Test case: success", func(t *testing.T) {
		// Mock InsertOne return value
		mockInsertResult := &mongo.InsertOneResult{
			InsertedID: primitive.NewObjectID(),
		}
		mockCollection = new(MockMongoCollection)
		// Mock logic cho FindOne - ensure this returns *MockSingleResult
		mockSingleResult := new(MockSingleResult)
		mockSingleResult.On("Decode", mock.Anything).Return(nil) // Mock no existing user
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockSingleResult)

		// Mock logic cho InsertOne
		mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(mockInsertResult, nil)

		// Tạo request
		registerRequest := struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
			Profile  struct {
				PhoneNumber string `json:"phone_number"`
			} `json:"profile"`
		}{
			Username: "testuser",
			Password: "validPassword1!",
			Email:    "test@example.com",
			Profile: struct {
				PhoneNumber string `json:"phone_number"`
			}{
				PhoneNumber: "+84345678904",
			},
		}
		body, _ := json.Marshal(registerRequest)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		// Gửi request và kiểm tra kết quả
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "User registered successfully")

		// Kiểm tra expectation của mock
		mockCollection.AssertExpectations(t)
		mockSingleResult.AssertExpectations(t)
	})

    // Test case: Kiểm tra đã có người dùng tồn tại
    t.Run("Test case: user exists", func(t *testing.T) {
        // Mock data cho test case người dùng đã tồn tại
        existingUser := models.User{
            Username: "testuser",
            Email:    "test@example.com",
            Password: "validPassword1!",
        }
		mockCollection = new(MockMongoCollection)
        // Mock logic cho FindOne
        mockResult := new(MockSingleResult)
        mockResult.On("Decode", mock.Anything).Return(existingUser, nil) // Người dùng đã tồn tại
        mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockResult)

        // Tạo request
        registerRequest := struct {
            Username string `json:"username"`
            Password string `json:"password"`
            Email    string `json:"email"`
        }{
            Username: "testuser",
            Password: "validPassword1!",
            Email:    "test@example.com",
        }
        body, _ := json.Marshal(registerRequest)
        req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
        w := httptest.NewRecorder()

        // Gửi request và kiểm tra kết quả
        r.ServeHTTP(w, req)

        assert.Equal(t, http.StatusConflict, w.Code)
        assert.Contains(t, w.Body.String(), "Email or username or phone already exists.")

        // Kiểm tra expectation của mock
        mockResult.AssertExpectations(t)
        mockCollection.AssertExpectations(t)
    })
}

func TestLogout(t *testing.T) {
    gin.SetMode(gin.TestMode)

    r := gin.Default()
    //mockCollection := new(MockMongoCollection)

    r.POST("/logout", func(c *gin.Context) {
        // Lấy Access Token từ cookie
		_, err := c.Cookie("accessToken")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Access Token not provided",
			})
			return
		}

		// Lấy Refresh Token từ cookie
		_, err = c.Cookie("refreshToken")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Refresh Token not provided",
			})
			return
		}

		// // Giả lập việc xác thực Access Token
        // accessClaims, err := mockMiddleware.VerifyJWT(accessToken, true)
        // if err == nil {
        //     if stdClaims, ok := accessClaims.(*jwt.StandardClaims); ok {
        //         middlewares.BlacklistedTokens[accessToken] = stdClaims.ExpiresAt
        //     }
        // }

        // // Giả lập việc xác thực Refresh Token
        // refreshClaims, err := mockMiddleware.VerifyJWT(refreshToken, false)
        // if err == nil {
        //     if stdClaims, ok := refreshClaims.(*jwt.StandardClaims); ok {
        //         middlewares.BlacklistedTokens[refreshToken] = stdClaims.ExpiresAt
        //     }
        // }

		// Gọi hàm reset cookie để xóa cookies
		// err = resetAuthCookies(c)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		// Trả về thông báo thành công
		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
    })

    // Test case: Không có Access Token trong cookies
	t.Run("Test case: missing access token", func(t *testing.T) {
		// Giả lập request thiếu Access Token cookie
		req, _ := http.NewRequest(http.MethodPost, "/logout", nil)
		req.AddCookie(&http.Cookie{Name: "refreshToken", Value: "valid-refresh-token"})

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Kiểm tra mã phản hồi và lỗi
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Access Token not provided"}`, w.Body.String())
	})

    // Test case: Không có Refresh Token trong cookies
	t.Run("Test case: missing refresh token", func(t *testing.T) {
		// Giả lập request thiếu Refresh Token cookie
		req, _ := http.NewRequest(http.MethodPost, "/logout", nil)
		req.AddCookie(&http.Cookie{Name: "accessToken", Value: "valid-access-token"})

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Kiểm tra mã phản hồi và lỗi
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Refresh Token not provided"}`, w.Body.String())
	})
	// Test case: Kiểm tra khi logout thành công
    t.Run("Test case: Success", func(t *testing.T) {
        accessToken := "valid_access_token"
        refreshToken := "valid_refresh_token"

        // Tạo request với cookies
        req, _ := http.NewRequest("POST", "/logout", nil)
        req.AddCookie(&http.Cookie{Name: "accessToken", Value: accessToken})
        req.AddCookie(&http.Cookie{Name: "refreshToken", Value: refreshToken})
        
        w := httptest.NewRecorder()

        r.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
        assert.Contains(t, w.Body.String(), "Logout successful")
    })
}

func TestForgotPassword(t *testing.T) {
    gin.SetMode(gin.TestMode)

    r := gin.Default()
    mockCollection := new(MockMongoCollection)

    r.POST("/forgot-password", func(c *gin.Context) {
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
	
		// Sử dụng mockCollection thay vì kết nối thật tới MongoDB
		result := mockCollection.FindOne(context.TODO(), map[string]interface{}{
			"email": request.Email,
		})

		var user models.User
		if err := result.Decode(&user); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found with this email"})
			return
		}
	
	
		// Tạo token ngẫu nhiên và hash nó để lưu vào DB
		// Lưu token vào cơ sở dữ liệu
		// Gửi email
		// Gọi hàm gửi email
	
		c.JSON(http.StatusOK, gin.H{
			"message": "Password reset link sent to your email",
		})
	})

	// Test case 1: Success
	t.Run("Test case: Success", func(t *testing.T) {
		// Mock data cho test case success
		userID := primitive.NewObjectID()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("dath2024@"), bcrypt.DefaultCost)
		validUser := models.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
			Password: string(hashedPassword),
			IsActive: true,
			Role:     "admin",
		}
		mockResult := new(MockSingleResult)
		// Mock logic cho FindOne và Decode
		mockResult.On("Decode", mock.Anything).Return(validUser, nil)
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockResult)

		// Create request body
		request := struct {
			Email string `json:"email"`
		}{
			Email: "test@example.com",
		}
		body, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", "/forgot-password", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	
		// Check response code and message
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Password reset link sent to your email")
	})
	
    // Test case 2: Missing email in request
	t.Run("Test case: Missing email", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/forgot-password", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid request format"}`, w.Body.String())
	})
	// Test case 3: User not found with email
	t.Run("Test case: User not found", func(t *testing.T) {
		mockCollection = new(MockMongoCollection)
		mockResult := new(MockSingleResult)
	
		// Mock logic cho user not found
		mockResult.On("Decode", mock.Anything).Return(models.User{}, assert.AnError)
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockResult)
	
		// Create request body
		request := struct {
			Email string `json:"email"`
		}{
			Email: "test@gmail.com",
		}
		body, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", "/forgot-password", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	
		// Check response code and message
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "User not found with this email")
	})
}

type MockUser struct {
    ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Username               string             `bson:"username" json:"username"`
    Email                  string             `bson:"email" json:"email"`
    Password               string             `bson:"password" json:"password"`
    Role                   string             `bson:"role" json:"role"`
    IsActive               bool               `bson:"is_active" json:"is_active"`
    ResetPasswordToken     string             `bson:"reset_password_token,omitempty" json:"reset_password_token,omitempty"`
    ResetPasswordExpires   time.Time          `bson:"reset_password_expires,omitempty" json:"reset_password_expires,omitempty"`
}

func TestResetPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	mockCollection := new(MockMongoCollection)

	r.POST("/reset-password", func(c *gin.Context) {
		var request struct {
			Token       string `json:"token"`
			NewPassword string `json:"new_password"`
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
		var user models.User
		err := mockCollection.FindOne(context.TODO(), bson.M{
			"reset_password_token": hashedToken,
			"reset_password_expires": bson.M{
				"$gt": time.Now(), // Chỉ chấp nhận token còn hiệu lực
			},
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
		_, err = mockCollection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update password",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Password changed successfully",
		})
	})

	// Test case 1: Success
	t.Run("Test case: Success", func(t *testing.T) {
		// Mock data cho test case success
		userID := primitive.NewObjectID()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("dath2024@"), bcrypt.DefaultCost)
		validUser := models.User{
			ID:                   userID,
			Username:             "testuser",
			Email:                "test@example.com",
			Password:             string(hashedPassword),
			IsActive:             true,
			Role:                 "admin",
			//ResetPasswordToken:   utils.HashString("tokenreset"),
			//ResetPasswordExpires: time.Now().Add(1 * time.Hour), // Token còn hiệu lực
		}

		// Mock kết quả FindOne
		mockResult := new(MockSingleResult)
		mockResult.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
			user := args.Get(0).(*models.User)
			*user = validUser // Trả về user giả lập
		}).Return(nil)

		// Mock FindOne
		mockCollection.On("FindOne", mock.Anything, mock.MatchedBy(func(filter interface{}) bool {
			f, ok := filter.(bson.M)
			if !ok {
				return false
			}

			// Kiểm tra token và thời hạn trong filter
			tokenMatch := f["reset_password_token"] == utils.HashString("tokenreset")
			expiresMatch := f["reset_password_expires"].(bson.M)["$gt"].(time.Time).After(time.Now())

			return tokenMatch && expiresMatch
		})).Return(mockResult)

		// Mock UpdateOne
		mockCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{
			MatchedCount:  1,
			ModifiedCount: 1,
		}, nil)

		// Tạo request body
		request := struct {
			Token       string `json:"token"`
			NewPassword string `json:"new_password"`
		}{
			Token:       "tokenreset",
			NewPassword: "@dath2024",
		}
		body, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", "/reset-password", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Kiểm tra response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Password changed successfully")
	})
}
