package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dath-241/coin-price-be-go/services/admin_service/controllers"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUserByAdmin(t *testing.T) {
	t.Run("Test Get User By Admin", func(t *testing.T) {
		objID, err := primitive.ObjectIDFromHex("6488e1c4b5d1e40b2c93f3a0")
		if err != nil {
			// Xử lý lỗi ở đây
		}
		mockRepo := &repository.MockUserRepository{
			Users: map[string]interface{}{
				"6488e1c4b5d1e40b2c93f3a0": models.User{
					ID:       objID,
					Username: "test_user",
					Email:    "test_user@example.com",
					IsActive: true,
				},
			},
			Err: nil, // Không có lỗi trong test case này
		}
		

		r := gin.Default()
		r.GET("/users/:user_id", controllers.GetUserByAdmin(mockRepo))

		req, _ := http.NewRequest(http.MethodGet, "/users/6488e1c4b5d1e40b2c93f3a0", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test_user", response["username"])
		assert.Equal(t, "test_user@example.com", response["email"])
		assert.True(t, response["is_active"].(bool))
	})
}

func TestDeleteUserByAdmin(t *testing.T) {
	t.Run("Test Delete User By Admin", func(t *testing.T) {
		mockRepo := &repository.MockUserRepository{
			Users: map[string]interface{}{
				"6488e1c4b5d1e40b2c93f3a0": map[string]interface{}{
					"_id":      "6488e1c4b5d1e40b2c93f3a0", // ID của người dùng
					"username": "test_user",               // Tên người dùng
					"email":    "test_user@example.com",   // Email
					"is_active": true,                     // Trạng thái người dùng
				},
			},
			Err: nil, // Không có lỗi trong test case này
		}
		

		r := gin.Default()
		r.DELETE("/users/:user_id", controllers.DeleteUserByAdmin(mockRepo))

		// Test trường hợp xóa thành công
		req, _ := http.NewRequest(http.MethodDelete, "/users/6488e1c4b5d1e40b2c93f3a0", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User deleted successfully", response["message"])

		// Test trường hợp không tìm thấy người dùng
		req, _ = http.NewRequest(http.MethodDelete, "/users/6488e1c4b5d1e40b2c93f3a2", nil)
		rec = httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		var responseNotFound map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &responseNotFound)
		assert.NoError(t, err)
		assert.Equal(t, "User not found", responseNotFound["error"])
	})
}

func TestGetPaymentHistoryForAdmin(t *testing.T) {
	t.Run("Test Get Payment History For Admin", func(t *testing.T) {
		// Mock repository với dữ liệu thanh toán giả lập
		mockRepo := &repository.MockPaymentRepository{
			Payments: map[string]interface{}{
				"order1": models.Order{
					ID:                primitive.NewObjectID(),
					UserID:            primitive.NewObjectID(),
					OrderID:           "order1",
					OrderInfo:         "Product A",
					TransactionStatus: "Completed",
					Amount:            100,
				},
				"order2": models.Order{
					ID:                primitive.NewObjectID(),
					UserID:            primitive.NewObjectID(),
					OrderID:           "order2",
					OrderInfo:         "Product B",
					TransactionStatus: "Failed",
					Amount:            50,
				},
			},
			Err: nil, // Không có lỗi trong test case này
		}

		// Thiết lập router với mock repository
		r := gin.Default()
		r.GET("/payment-history", controllers.GetPaymentHistoryForAdmin(mockRepo))

		// Test trường hợp hợp lệ (có thanh toán)
		t.Run("Valid Payment History", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/payment-history", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			// Kiểm tra mã trạng thái HTTP
			assert.Equal(t, http.StatusOK, rec.Code)

			// Giải mã nội dung phản hồi
			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Kiểm tra cấu trúc và nội dung của phản hồi
			paymentHistory, ok := response["payment_history"].([]interface{})
			assert.True(t, ok)
			assert.Equal(t, len(mockRepo.Payments), len(paymentHistory))
		})

		// Test trường hợp không có thanh toán
		t.Run("No Payment History", func(t *testing.T) {
			emptyRepo := &repository.MockPaymentRepository{
				Payments: map[string]interface{}{},
			}
			r := gin.Default()
			r.GET("/payment-history", controllers.GetPaymentHistoryForAdmin(emptyRepo))

			req := httptest.NewRequest(http.MethodGet, "/payment-history", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			// Kiểm tra mã trạng thái HTTP
			assert.Equal(t, http.StatusOK, rec.Code)

			// Kiểm tra nội dung phản hồi
			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "No payment history found", response["message"])
		})
	})
}

func TestBanAccount(t *testing.T) {
	t.Run("Test Ban Account", func(t *testing.T) {
		// Mock user repository
		userID, _ := primitive.ObjectIDFromHex("6488e1c4b5d1e40b2c93f3a0")

		mockRepo := &repository.MockUserRepository{
			Users: map[string]interface{}{
				userID.Hex(): models.User{
					ID:       userID,
					IsActive: true,
				},
			},
		}

		// Set up Gin router and handler
		r := gin.Default()
		r.PUT("/ban-account/:user_id", controllers.BanAccount(mockRepo))

		// Case 1: Successful account ban
		t.Run("Ban Account Successfully", func(t *testing.T) {
			// Request to ban account
			req := httptest.NewRequest(http.MethodPut, "/ban-account/"+userID.Hex(), nil)
			rec := httptest.NewRecorder()

			// Call the handler
			r.ServeHTTP(rec, req)

			// Assert the response
			assert.Equal(t, http.StatusOK, rec.Code)

			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Account has been banned successfully", response["message"])

			// Verify that the user is banned
			updatedUser := mockRepo.Users[userID.Hex()].(models.User)
			assert.False(t, updatedUser.IsActive)
		})

		// Case 2: User not found
		t.Run("User Not Found", func(t *testing.T) {
			// Update mock repository to simulate no user found
			mockRepo.Users = make(map[string]interface{})

			// Request to ban non-existing account
			req := httptest.NewRequest(http.MethodPut, "/ban-account/"+userID.Hex(), nil)
			rec := httptest.NewRecorder()

			// Call the handler
			r.ServeHTTP(rec, req)

			// Assert the response
			assert.Equal(t, http.StatusNotFound, rec.Code)

			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "User not found", response["error"])
		})
	})
}

func TestActiveAccount(t *testing.T) {
	t.Run("Test Active Account", func(t *testing.T) {
		// Mock user repository
		userID, _ := primitive.ObjectIDFromHex("6488e1c4b5d1e40b2c93f3a0")

		mockRepo := &repository.MockUserRepository{
			Users: map[string]interface{}{
				userID.Hex(): models.User{
					ID:       userID,
					IsActive: false,
				},
			},
		}

		// Set up Gin router and handler
		r := gin.Default()
		r.PUT("/ban-account/:user_id", controllers.ActiveAccount(mockRepo))

		// Case 1: Successful account ban
		t.Run("Ban Account Successfully", func(t *testing.T) {
			// Request to ban account
			req := httptest.NewRequest(http.MethodPut, "/ban-account/"+userID.Hex(), nil)
			rec := httptest.NewRecorder()

			// Call the handler
			r.ServeHTTP(rec, req)

			// Assert the response
			assert.Equal(t, http.StatusOK, rec.Code)

			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Account has been active successfully", response["message"])

			// Verify that the user is banned
			updatedUser := mockRepo.Users[userID.Hex()].(models.User)
			assert.True(t, updatedUser.IsActive)
		})

		// Case 2: User not found
		t.Run("User Not Found", func(t *testing.T) {
			// Update mock repository to simulate no user found
			mockRepo.Users = make(map[string]interface{})

			// Request to ban non-existing account
			req := httptest.NewRequest(http.MethodPut, "/ban-account/"+userID.Hex(), nil)
			rec := httptest.NewRecorder()

			// Call the handler
			r.ServeHTTP(rec, req)

			// Assert the response
			assert.Equal(t, http.StatusNotFound, rec.Code)

			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "User not found", response["error"])
		})
	})
}

