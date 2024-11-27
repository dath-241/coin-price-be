package controllers_test

import (
	"bytes"
	"os"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dath-241/coin-price-be-go/services/admin_service/controllers"
	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

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
				Username: "test-user",
				Email:    "test_user@example.com",
				Profile: models.Profile{
					FullName:    "Test User",
					PhoneNumber: "8423456789",
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
			"username": "updated-user",
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
		assert.Equal(t, "updated-user", updatedUser.(models.User).Username)
		//assert.Equal(t, "987654321", updatedUser.(models.User).Profile.PhoneNumber)
	})

	// Case 2: Missing Authorization header
	t.Run("Missing Authorization", func(t *testing.T) {
		updateRequest := map[string]interface{}{
			"username": "update-user", // Trying to update with an existing username
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
			"username": "update-user", // Trying to update with an existing username
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
