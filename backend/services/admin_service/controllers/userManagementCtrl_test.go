package controllers_test

import (
	//"context"
	"bytes"
	"os"
	//"context"
	"encoding/json"
	//"fmt"
	//"time"

	//"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	//"time"

	"github.com/dath-241/coin-price-be-go/services/admin_service/controllers"
	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func TestGetCurrentUserInfo(t *testing.T) {
	userID, _ := primitive.ObjectIDFromHex("6488e1c4b5d1e40b2c93f3a0")
	mockRepo := &repository.MockUserRepository{
		Users: map[string]interface{}{
			userID.Hex(): models.User{
				ID: userID,
				Username: "testuser",
				Email: "testuser@example.com",
				Profile: models.Profile{
					FullName:    "Test User",
					PhoneNumber: "123456789",
					AvatarURL:   "https://example.com/avatar.jpg",
					Bio:         "Hello, I'm a test user.",
					DateOfBirth: "2000-01-01",
				},
				Role: "VIP-1",
			},
		},
	}

	r := gin.Default()
	r.GET("/current-user", controllers.GetCurrentUserInfo(mockRepo))

	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DB_NAME", "test_db")
	os.Setenv("ACCESS_TOKEN_TTL", "3600")
	
	// Tạo token hợp lệ
	validToken, err := middlewares.GenerateAccessToken("6488e1c4b5d1e40b2c93f3a0","VIP-1")
	if err != nil {
		t.Fatalf("failed to generate valid token: %v", err)
	}

	defer os.Unsetenv("MONGO_URI")
	defer os.Unsetenv("MONGO_DB_NAME")
	defer os.Unsetenv("ACCESS_TOKEN_TTL") // Dọn dẹp biến môi trường

	// Case 1: Successful retrieval of user info
	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/current-user", nil)
		req.Header.Set("Authorization", validToken)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", response["username"])
		assert.Equal(t, "Test User", response["name"])
		assert.Equal(t, "123456789", response["phone_number"])
	})

	// Case 2: Missing Authorization header
	t.Run("MissingAuthorization", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/current-user", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authorization header required", response["error"])
	})

	// Case 3: Invalid Token
	t.Run("InvalidToken", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/current-user", nil)
		req.Header.Set("Authorization", "invalid_token")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid token")
	})
}


func TestUpdateUserProfile(t *testing.T) {
	// Mock user repository
	userID, _ := primitive.ObjectIDFromHex("6488e1c4b5d1e40b2c93f3a0")
	mockRepo := &repository.MockUserRepository{
		Users: map[string]interface{}{
			userID.Hex(): models.User{
				ID: userID,
				Username: "test_user",
				Email:    "test_user@example.com",
				Profile: models.Profile{
					FullName:    "Test User",
					PhoneNumber: "123456789",
					AvatarURL:   "https://example.com/avatar.jpg",
					Bio:         "Hello, I'm a test user.",
					DateOfBirth: "2000-01-01",
				},
				Role: "VIP-1",
			},
		},
	}

	// Set up Gin router and handler
	r := gin.Default()
	r.PUT("/update", controllers.UpdateUserProfile(mockRepo))

	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DB_NAME", "test_db")
	os.Setenv("ACCESS_TOKEN_TTL", "3600")

	// Generate a valid token for the user
	validToken, err := middlewares.GenerateAccessToken(userID.Hex(), "VIP-1")
	if err != nil {
		t.Fatalf("failed to generate valid token: %v", err)
	}

	defer os.Unsetenv("MONGO_URI")
	defer os.Unsetenv("MONGO_DB_NAME")
	defer os.Unsetenv("ACCESS_TOKEN_TTL") // Dọn dẹp biến môi trường

	// Case 1: Valid token and successful update
	t.Run("Valid Update", func(t *testing.T) {
		updateRequest := map[string]interface{}{
			"username": "updated_user",
		}
		reqBody, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(reqBody))
		req.Header.Set("Authorization", validToken)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User information updated successfully", response["message"])

		// Verify that the user data was updated
		updatedUser, found := mockRepo.Users[userID.Hex()]
		assert.True(t, found)
		assert.Equal(t, "updated_user", updatedUser.(models.User).Username)
		//assert.Equal(t, "987654321", updatedUser.(models.User).Profile.PhoneNumber)
	})

	// Case 2: Missing Authorization header
	t.Run("Missing Authorization", func(t *testing.T) {
		updateRequest := map[string]interface{}{
			"username": "update_user", // Trying to update with an existing username
		}
		reqBody, _ := json.Marshal(updateRequest)
		req, _ := http.NewRequest(http.MethodPut, "/update", bytes.NewReader(reqBody))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authorization header required", response["error"])
	})

	// Case 3: Invalid Token
	t.Run("Invalid Token", func(t *testing.T) {
		updateRequest := map[string]interface{}{
			"username": "update_user", // Trying to update with an existing username
		}
		reqBody, _ := json.Marshal(updateRequest)
		req, _ := http.NewRequest(http.MethodPut, "/update", bytes.NewReader(reqBody))
		req.Header.Set("Authorization", "invalid_token")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid token")
	})

	// Case 4: Username already in use
	t.Run("Username Conflict", func(t *testing.T) {
		// Adding another user with a different ID
		otherUserID, _ := primitive.ObjectIDFromHex("6488e1c4b5d1e40b2c93f3b1")
		mockRepo.Users[otherUserID.Hex()] = models.User{
			ID:       otherUserID,
			Username: "existing_user",
			Email:    "existing_user@example.com",
		}

		updateRequest := map[string]interface{}{
			"username": "existing_user", // Trying to update with an existing username
		}
		reqBody, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(reqBody))
		req.Header.Set("Authorization", validToken)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Username already in use", response["error"])
	})
}

func TestChangePassword(t *testing.T) {
	// Mock user repository
	userID, _ := primitive.ObjectIDFromHex("6488e1c4b5d1e40b2c93f3a0")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("current_password"), bcrypt.DefaultCost)

	mockRepo := &repository.MockUserRepository{
		Users: map[string]interface{}{
			userID.Hex(): models.User{
				ID:       userID,
				Password: string(hashedPassword),
				Role:     "VIP-1",
			},
		},
	}

	// Set up Gin router and handler
	r := gin.Default()
	r.POST("/change-password", controllers.ChangePassword(mockRepo))

	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DB_NAME", "test_db")
	os.Setenv("ACCESS_TOKEN_TTL", "3600")

	// Generate a valid token for the user
	validToken, err := middlewares.GenerateAccessToken(userID.Hex(), "VIP-1")
	if err != nil {
		t.Fatalf("failed to generate valid token: %v", err)
	}

	defer os.Unsetenv("MONGO_URI")
	defer os.Unsetenv("MONGO_DB_NAME")
	defer os.Unsetenv("ACCESS_TOKEN_TTL") // Dọn dẹp biến môi trường

	// Case 1: Valid token and successful password change
	t.Run("Valid Password Change", func(t *testing.T) {
		changeRequest := map[string]interface{}{
			"current_password": "current_password",
			"new_password":     "NewPassword@123",
		}
		reqBody, _ := json.Marshal(changeRequest)
		req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(reqBody))
		req.Header.Set("Authorization", validToken)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Password updated successfully", response["message"])

		// Verify that the password was updated
		updatedUser := mockRepo.Users[userID.Hex()].(models.User)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("NewPassword@123")))
	})

	// Case 2: Missing Authorization header
	t.Run("Missing Authorization", func(t *testing.T) {
		changeRequest := map[string]interface{}{
			"current_password": "current_password",
			"new_password":     "NewPassword@123",
		}
		reqBody, _ := json.Marshal(changeRequest)
		req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(reqBody))
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authorization header required", response["error"])
	})

	// Case 3: Invalid current password
	t.Run("Invalid Current Password", func(t *testing.T) {
		changeRequest := map[string]interface{}{
			"current_password": "wrong_password",
			"new_password":     "NewPassword@123",
		}
		reqBody, _ := json.Marshal(changeRequest)
		req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(reqBody))
		req.Header.Set("Authorization", validToken)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Current password is incorrect", response["error"])
	})
	// 	err := json.Unmarshal(rec.Body.Bytes(), &response)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, "Password must contain at least 8 characters, including letters, numbers, and special characters.", response["error"])
	// })
}

func TestChangeEmail(t *testing.T) {
	// Mock user repository
	userID, _ := primitive.ObjectIDFromHex("6488e1c4b5d1e40b2c93f3a0")

	mockRepo := &repository.MockUserRepository{
		Users: map[string]interface{}{
			userID.Hex(): models.User{
				ID:    userID,
				Email: "currentemail@example.com",
				Role:  "VIP-1",
			},
		},
	}

	// Set up Gin router and handler
	r := gin.Default()
	r.PUT("/change-email", controllers.ChangeEmail(mockRepo))

	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DB_NAME", "test_db")
	os.Setenv("ACCESS_TOKEN_TTL", "3600")

	// Generate a valid token for the user
	validToken, err := middlewares.GenerateAccessToken(userID.Hex(), "VIP-1")
	if err != nil {
		t.Fatalf("failed to generate valid token: %v", err)
	}

	defer os.Unsetenv("MONGO_URI")
	defer os.Unsetenv("MONGO_DB_NAME")
	defer os.Unsetenv("ACCESS_TOKEN_TTL") // Dọn dẹp biến môi trường

	// Case 1: Valid email update
	t.Run("Valid Email Update", func(t *testing.T) {
		changeRequest := map[string]interface{}{
			"email": "newemail@example.com",
		}
		reqBody, _ := json.Marshal(changeRequest)
		req := httptest.NewRequest(http.MethodPut, "/change-email", bytes.NewReader(reqBody))
		req.Header.Set("Authorization", validToken)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Email updated successfully", response["message"])

		// Verify that the email was updated
		updatedUser := mockRepo.Users[userID.Hex()].(models.User)
		assert.Equal(t, "newemail@example.com", updatedUser.Email)
	})
}

func TestDeleteCurrentUser(t *testing.T) {
	// Mock user repository
	userID, _ := primitive.ObjectIDFromHex("6488e1c4b5d1e40b2c93f3a0")

	mockRepo := &repository.MockUserRepository{
		Users: map[string]interface{}{
			userID.Hex(): models.User{
				ID: userID,
			},
		},
	}

	// Set up Gin router and handler
	r := gin.Default()
	r.PUT("/delete-account", controllers.DeleteCurrentUser(mockRepo))

	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DB_NAME", "test_db")
	os.Setenv("ACCESS_TOKEN_TTL", "3600")

	// Generate a valid token for the user
	validToken, err := middlewares.GenerateAccessToken(userID.Hex(), "VIP-1")
	if err != nil {
		t.Fatalf("failed to generate valid token: %v", err)
	}

	defer os.Unsetenv("MONGO_URI")
	defer os.Unsetenv("MONGO_DB_NAME")
	defer os.Unsetenv("ACCESS_TOKEN_TTL") // Dọn dẹp biến môi trường

	// Case 1: Valid user deletion
	t.Run("Valid User Deletion", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]interface{}{})
		req := httptest.NewRequest(http.MethodPut, "/delete-account", bytes.NewReader(reqBody))
		req.Header.Set("Authorization", validToken)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User account deleted successfully", response["message"])

		// Verify that the user is deleted
		_, found := mockRepo.Users[userID.Hex()]
		assert.False(t, found)
	})

	// Case 2: User not found (deleted or never existed)
	t.Run("User Not Found", func(t *testing.T) {
		// Delete the user from the mock repository
		delete(mockRepo.Users, userID.Hex())

		reqBody, _ := json.Marshal(map[string]interface{}{})
		req := httptest.NewRequest(http.MethodPut, "/delete-account", bytes.NewReader(reqBody))
		req.Header.Set("Authorization", validToken)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User not found", response["error"])
	})

	// Case 3: Invalid token
	t.Run("Invalid Token", func(t *testing.T) {
		invalidToken := "invalidToken"
		reqBody, _ := json.Marshal(map[string]interface{}{})
		req := httptest.NewRequest(http.MethodPut, "/delete-account", bytes.NewReader(reqBody))
		req.Header.Set("Authorization", invalidToken)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid token", response["error"])
	})
}


// func TestGetPaymentHistory(t *testing.T) {
// 	// Tạo user_id giả
// 	userID := primitive.NewObjectID()

// 	// Mock repository với một số payment data
// 	mockRepo := &repository.MockPaymentRepository{
// 		Payments: map[string]interface{}{
// 			userID.Hex(): models.Order{
// 				OrderID:          "order1",
// 				UserID:           userID,
// 				OrderInfo:        "Order details",
// 				TransactionStatus: "Completed",
// 				Amount:            100,
// 			},
// 		},
// 	}

// 	// Set up Gin router and handler
// 	r := gin.Default()
// 	r.GET("/payment-history", controllers.GetPaymentHistory(mockRepo))

// 	// Generate a valid token for the user using the middlewares.GenerateAccessToken
// 	validToken, err := middlewares.GenerateAccessToken(userID.Hex(), "VIP-1")
// 	if err != nil {
// 		t.Fatalf("failed to generate valid token: %v", err)
// 	}

// 	// Case 1: Valid payment history
// 	t.Run("Valid Payment History", func(t *testing.T) {
// 		req, _ := http.NewRequest(http.MethodGet, "/payment-history", nil)
// 		req.Header.Set("Authorization", validToken)
// 		rec := httptest.NewRecorder()
	
// 		// Gọi handler gin
// 		r.ServeHTTP(rec, req)
	
// 		// Kiểm tra code trả về
// 		assert.Equal(t, http.StatusOK, rec.Code)
	
// 		var response map[string]interface{}
// 		err = json.Unmarshal(rec.Body.Bytes(), &response)
// 		assert.NoError(t, err)
	
// 		// Kiểm tra xem "payment_history" có tồn tại không
// 		paymentHistory, ok := response["payment_history"].([]interface{})
// 		if !ok {
// 			t.Fatalf("Expected payment_history to be []interface{}, but got %T", response["payment_history"])
// 		}
	
// 		// Kiểm tra xem paymentHistory có phần tử nào không
// 		if len(paymentHistory) == 0 {
// 			t.Errorf("No payment history found, response is empty")
// 		} else {
// 			// Kiểm tra phần tử trong paymentHistory
// 			paymentDetails, ok := paymentHistory[0].(map[string]interface{})
// 			if !ok {
// 				t.Fatalf("Expected paymentDetails to be map[string]interface{}, but got %T", paymentHistory[0])
// 			}
	
// 			// Kiểm tra các trường trong payment
// 			assert.Equal(t, "order1", paymentDetails["order_id"])
// 			assert.Equal(t, userID.Hex(), paymentDetails["user_id"])
// 			assert.Equal(t, "Completed", paymentDetails["transaction_status"])
// 			assert.Equal(t, float64(100), paymentDetails["amount"])
// 		}
// 	})

// 	// Case 2: No payment history found
// 	t.Run("No Payment History Found", func(t *testing.T) {
// 		// Mock repository không có payment
// 		mockRepo.Payments = map[string]interface{}{}

// 		req, _ := http.NewRequest(http.MethodGet, "/payment-history", nil)
// 		req.Header.Set("Authorization",validToken)
// 		rec := httptest.NewRecorder()

// 		r.ServeHTTP(rec, req)

// 		assert.Equal(t, http.StatusOK, rec.Code)

// 		var response map[string]interface{}
// 		err := json.Unmarshal(rec.Body.Bytes(), &response)
// 		assert.NoError(t, err)

// 		// Kiểm tra thông báo không có lịch sử thanh toán
// 		assert.Equal(t, "No payment history found", response["message"])
// 	})

// 	// Case 3: Repository error
// 	t.Run("Repository Error", func(t *testing.T) {
// 		// Simulate a repository error
// 		mockRepo.Err = fmt.Errorf("repository error")

// 		req, _ := http.NewRequest(http.MethodGet, "/payment-history", nil)
// 		req.Header.Set("Authorization",validToken)
// 		rec := httptest.NewRecorder()

// 		r.ServeHTTP(rec, req)

// 		assert.Equal(t, http.StatusInternalServerError, rec.Code)

// 		var response map[string]interface{}
// 		err := json.Unmarshal(rec.Body.Bytes(), &response)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "Error fetching payment history", response["error"])
// 	})
// }
