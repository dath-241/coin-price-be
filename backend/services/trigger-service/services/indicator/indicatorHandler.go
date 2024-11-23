package services

import (
	"net/http"

	"github.com/dath-241/coin-price-be-go/services/admin_service/config"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SetAdvancedIndicatorAlert creates a new advanced indicator alert
// @Summary Create an advanced indicator alert
// @Description Create an alert with the given indicator type and settings
// @Tags Indicators
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param body body models.Indicator true "Indicator alert details"
// @Success 201 {object} models.ResponseIndicatorCreated "Indicator alert created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body or invalid indicator type"
// @Failure 500 {object} models.ErrorResponse "Failed to create indicator alert"
// @Router /api/v1/vip3/indicators [post]
// PostIndicator creates a new advanced indicator alert
func SetAdvancedIndicatorAlert(c *gin.Context) {
	var newIndicator models.Indicator

	if err := c.ShouldBindJSON(&newIndicator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if newIndicator.Indicator != "EMA" && newIndicator.Indicator != "BollingerBands" && newIndicator.Indicator != "Custom" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid indicator type"})
		return
	}

	newIndicator.ID = primitive.NewObjectID().Hex()

	// Insert into the IndicatorCollection
	_, err := config.IndicatorCollection.InsertOne(c, newIndicator)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create indicator alert"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Indicator created successfully",
		"alert_id": newIndicator.ID,
	})
}
