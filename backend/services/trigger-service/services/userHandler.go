package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	config "github.com/dath-241/coin-price-be-go/services/admin_service/config"
	models "github.com/dath-241/coin-price-be-go/services/trigger-service/models"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/repositories"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUser handles the creation of a new user.
// @Summary Create a user
// @Description Create a new user with the given details
// @Tags Users
// @Accept json
// @Produce json
// @Param body body models.User true "User details"
// @Success 201 {object} models.ResponseUserCreated "User created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body or missing email"
// @Failure 500 {object} models.ErrorResponse "Failed to create user"
// @Router /api/v1/users [post]
func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if newUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// Create a new unique ID for the user
	newUser.ID = primitive.NewObjectID().Hex()
	newUser.Alerts = []models.Alert{}

	// Set a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert the new user into the database
	_, err := config.AlertCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return a success response with the new user ID
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user_id": newUser.ID,
	})
}

// GetUserAlerts retrieves all alerts for a user by their ID.
// GetUserAlerts retrieves all alerts for a user by their ID.
// @Summary Get user alerts
// @Description Retrieve all alerts for a user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.ResponseUserAlerts "List of user alerts"
// @Failure 500 {object} models.ErrorResponse "Failed to retrieve alerts"
// @Router /api/v1/users/{id}/alerts [get]
func GetUserAlerts(c *gin.Context) {
	// Get user ID from URL parameter
	userID := c.Param("id")

	// Fetch user alerts from the database
	alerts, err := repositories.GetUserAlerts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the alerts to the client
	c.JSON(http.StatusOK, alerts)
}

// NotifyUser sends a notification email for a user's alerts.
// @Summary Notify user of alerts
// @Description Send a notification email to the user for their alerts
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.ResponseNotificationSent "Notification sent successfully"
// @Failure 500 {object} models.ErrorResponse "Failed to send notification"
// @Router /api/v1/users/{id}/alerts/notify [post]
func NotifyUser(c *gin.Context) {
	userID := c.Param("id")
	err := NotifyUserTriggers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Notification sent"})
}

func NotifyUserTriggers(userID string) error {

	// Retrieve user details from the repository
	user, err := repositories.GetUserByID(userID)
	log.Println(user)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Check if the user has an email
	if user.Email == "" {
		return fmt.Errorf("user email is missing")
	}

	// Retrieve user alerts
	alerts, err := repositories.GetUserAlerts(userID)
	log.Println("Email sent successfully")
	if err != nil {
		return fmt.Errorf("failed to retrieve alerts")
	}

	// Prepare the email subject and HTML body
	subject := "Your Trigger Alerts"
	htmlBody := "<h1>Trigger Alerts</h1><ul>"
	for _, alert := range alerts {
		htmlBody += fmt.Sprintf("<li><strong>%s:</strong> %s</li>", alert.Symbol, alert.Message)
	}
	htmlBody += "</ul>"

	// Send the email notification to the user

	if err := utils.SendAlertEmail(user.Email, subject, htmlBody); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
