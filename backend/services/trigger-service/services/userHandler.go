package services

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/dath-241/coin-price-be-go/services/trigger-service/models"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/repositories"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/utils"
	"github.com/gin-gonic/gin"
)

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

// NotifyUser sends a notification email for a user's alerts.
// @Summary Notify user of alerts
// @Description Send a notification email to the user for their alerts
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path string true "User ID"
// @Success 200 {object} models.ResponseNotificationSent "Notification sent successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Failed to send notification"
// @Security ApiKeyAuth
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
