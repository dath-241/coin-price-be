package services

import (
	"context"
	"log"
	"net/http"
	"time"

	modelsAD "github.com/dath-241/coin-price-be-go/services/admin_service/models"
	models "github.com/dath-241/coin-price-be-go/services/trigger-service/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	config "github.com/dath-241/coin-price-be-go/services/admin_service/config"
	middlewares "github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
)

// Handler to create an alert
// @Summary Create an alert
// @Description Create a new alert with the given details
// @Tags Alerts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param body body models.Alert true "Alert details"
// @Success 201 {object} models.ResponseAlertCreated "Successfully created alert"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Failed to create alert"
// @Security ApiKeyAuth
// @Router /api/v1/vip2/alerts [post]
func CreateAlert(c *gin.Context) {
	var newAlert models.Alert
	if err := c.ShouldBindJSON(&newAlert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Lấy token từ header Authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Xác thực token
	claims, err := middlewares.VerifyJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Lấy userID từ claims trong token
	currentUserID := claims.UserID
	log.Println(currentUserID)
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	// Chuyển user_id thành ObjectID
	objID, err := primitive.ObjectIDFromHex(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Kiểm tra số lượng alert hiện tại
	userCollection := config.DB.Collection("User")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user modelsAD.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if len(user.Alerts) >= 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum alert limit reached"})
		return
	}

	// Add default or new values for additional fields
	newAlert.ID = primitive.NewObjectID()
	newAlert.IsActive = true
	currentTime := primitive.NewDateTimeFromTime(time.Now())
	newAlert.CreatedAt = currentTime
	newAlert.UpdatedAt = currentTime

	if newAlert.MaxRepeatCount == 0 {
		newAlert.MaxRepeatCount = 5
	}
	newAlert.UserID = currentUserID

	if _, err := config.AlertCollection.InsertOne(ctx, newAlert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save alert"})
		return
	}

	// Lưu chỉ ID của alert vào user alerts
	user.Alerts = append(user.Alerts, newAlert.ID.Hex())
	// Cập nhật lại user trong cơ sở dữ liệu
	update := bson.M{
		"$set": bson.M{
			"alerts":     user.Alerts,
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	if _, err := userCollection.UpdateOne(ctx, bson.M{"_id": objID}, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user alerts"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Alert created successfully",
		"alert_id": newAlert.ID.Hex(),
	})
}

// Handler to retrieve all alerts
// @Summary Get all alerts
// @Description Retrieve all alerts, optionally filter by type
// @Tags Alerts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param type query string false "Filter by alert type (e.g., new_listing, delisting)"
// @Success 200 {array} models.ResponseAlertList "List of alerts"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Failed to retrieve alerts"
// @Security ApiKeyAuth
// @Router /api/v1/vip2/alerts [get]
func GetAlerts(c *gin.Context) {
	// Lấy token từ header Authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Xác thực token
	claims, err := middlewares.VerifyJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Lấy userID từ claims trong token
	currentUserID := claims.UserID
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	// Áp dụng filter tìm kiếm, chỉ lấy alert của user hiện tại
	alertType := c.Query("type")
	filter := bson.M{"user_id": currentUserID}
	if alertType != "" {
		filter["type"] = alertType
	}

	var results []models.Alert
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.AlertCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve alerts"})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse alerts"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// Handler to get an alert by ID
// @Summary Get an alert by ID
// @Description Retrieve an alert by its unique identifier
// @Tags Alerts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path string true "Alert ID"
// @Success 200 {object} models.ResponseAlertDetail "Alert details"
// @Failure 400 {object} models.ErrorResponse "Invalid alert ID"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Alert not found"
// @Security ApiKeyAuth
// @Router /api/v1/vip2/alerts/{id} [get]
func GetAlert(c *gin.Context) {
	// Lấy token từ header Authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Xác thực token
	claims, err := middlewares.VerifyJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Lấy userID từ claims trong token
	currentUserID := claims.UserID
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	// Lấy alert ID từ param
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	// Tìm alert và kiểm tra xem alert có thuộc về user hiện tại không
	var alert models.Alert
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = config.AlertCollection.FindOne(ctx, bson.M{"_id": objectId, "user_id": currentUserID}).Decode(&alert)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// Handler to delete an alert by ID
// @Summary Delete an alert
// @Description Delete an alert by its unique identifier
// @Tags Alerts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path string true "Alert ID"
// @Success 200 {object} models.ResponseAlertDeleted "Alert deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid alert ID"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Alert not found"
// @Security ApiKeyAuth
// @Router /api/v1/vip2/alerts/{id} [delete]
func DeleteAlert(c *gin.Context) {
	// Lấy token từ header Authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Xác thực token
	claims, err := middlewares.VerifyJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Lấy userID từ claims trong token
	currentUserID := claims.UserID
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	// Tìm alert và kiểm tra xem nó có thuộc về người dùng hiện tại không
	var alert models.Alert
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = config.AlertCollection.FindOne(ctx, bson.M{"_id": objectId, "user_id": currentUserID}).Decode(&alert)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	// Xoá alert từ AlertCollection
	result, err := config.AlertCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	// Cập nhật mảng alerts của user để loại bỏ alert đã xoá
	userCollection := config.DB.Collection("User")
	update := bson.M{
		"$pull": bson.M{"alerts": id}, // Xoá id khỏi mảng alerts
	}
	if _, err := userCollection.UpdateOne(ctx, bson.M{"_id": currentUserID}, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user alerts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alert deleted successfully"})
}

// Handler to retrieve new and delisted symbols
// @Summary Get new and delisted symbols
// @Description Retrieve new and delisted symbols from Binance
// @Tags Alerts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} models.ResponseNewDelistedSymbols "List of new and delisted symbols"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Failed to retrieve symbols"
// @Security ApiKeyAuth
// @Router /api/v1/vip2/symbols-alerts [get]
func GetSymbolAlerts(c *gin.Context) {
	// Lấy token từ header Authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Xác thực token
	claims, err := middlewares.VerifyJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Lấy userID từ claims trong token
	currentUserID := claims.UserID
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	newSymbols, delistedSymbols, err := FetchSymbolsFromBinance()
	if err != nil {
		log.Printf("Error fetching symbol data: %v", err)
		return
	}

	if newSymbols == nil {
		newSymbols = []string{}
	}
	if delistedSymbols == nil {
		delistedSymbols = []string{}
	}

	response := gin.H{
		"new_symbols":      newSymbols,
		"delisted_symbols": delistedSymbols,
	}

	c.JSON(http.StatusOK, response)
}

// Handler to set a symbol alert for new or delisted symbols
// Handler to set a symbol alert for new or delisted symbols
// @Summary Set an alert for new or delisted symbols
// @Description Set a new alert for symbols that are newly listed or delisted
// @Tags Alerts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param body body models.Alert true "Alert details"
// @Success 201 {object} models.ResponseSetSymbolAlert "Successfully created alert for symbol"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Failed to create alert for symbol"
// @Security ApiKeyAuth
// @Router /api/v1/vip2/alerts/symbol [post]
func SetSymbolAlert(c *gin.Context) {

	var newAlert models.Alert
	if err := c.ShouldBindJSON(&newAlert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// if (newAlert.Type != "new_listing" && newAlert.Type != "delisting") || newAlert.NotificationMethod == "" || len(newAlert.Symbols) == 0 || newAlert.Frequency == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid fields"})
	// 	return
	// }

	newAlert.ID = primitive.NewObjectID()
	newAlert.IsActive = true

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.AlertCollection.InsertOne(ctx, newAlert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create alert"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Alert created successfully",
		"alert_id": newAlert.ID.Hex(),
	})

}
