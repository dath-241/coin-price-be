package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	alert "github.com/dath-241/coin-price-be-go/services/trigger-service/models/alert"
	user "github.com/dath-241/coin-price-be-go/services/trigger-service/models/user"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/repositories"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

// CreateUser handles the creation of a new user.
func CreateUser(c *gin.Context) {
	var newUser user.User
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
	newUser.Alerts = []alert.Alert{}

	// Set a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert the new user into the database
	_, err := utils.AlertCollection.InsertOne(ctx, newUser)
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
    if err != nil {
        return fmt.Errorf("user not found")
    }

    // Check if the user has an email
    if user.Email == "" {
        return fmt.Errorf("user email is missing")
    }

    // Retrieve user alerts
    alerts, err := repositories.GetUserAlerts(userID)
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
