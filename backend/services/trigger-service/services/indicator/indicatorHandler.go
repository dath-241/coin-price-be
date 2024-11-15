package services

import (
	"net/http"

	"github.com/dath-241/coin-price-be-go/services/trigger-service/models/indicator"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostIndicator creates a new advanced indicator alert
func SetAdvancedIndicatorAlert(c *gin.Context) {
	var newIndicator indicator.Indicator

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
	_, err := utils.IndicatorCollection.InsertOne(c, newIndicator)
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
