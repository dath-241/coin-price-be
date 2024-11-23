package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

