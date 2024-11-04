package services

import (
    "context"
    "net/http"
    "time"

    "github.com/dath-241/coin-price-be-go/services/price-service/models/alert"
    "github.com/dath-241/coin-price-be-go/utils"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Create a new alert
func CreateAlert(c *gin.Context) {
    var newAlert models.Alert
    if err := c.ShouldBindJSON(&newAlert); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if newAlert.Condition == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid condition"})
        return
    }

    newAlert.ID = primitive.NewObjectID()
    newAlert.IsActive = true

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, err := utils.AlertCollection.InsertOne(ctx, newAlert)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create alert"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message":  "Alert created successfully",
        "alert_id": newAlert.ID.Hex(),
    })
}

// Get all alerts
func GetAlerts(c *gin.Context) {
    var results []models.Alert

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    cursor, err := utils.AlertCollection.Find(ctx, bson.M{})
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

// Get an alert by ID
func GetAlert(c *gin.Context) {
    id := c.Param("id")
    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
        return
    }

    var alert models.Alert
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    err = utils.AlertCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&alert)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
        return
    }

    c.JSON(http.StatusOK, alert)
}

// Delete an alert by ID
func DeleteAlert(c *gin.Context) {
    id := c.Param("id")
    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    filter := bson.M{"_id": objectId}
    result, err := utils.AlertCollection.DeleteOne(ctx, filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete alert"})
        return
    }

    if result.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Alert deleted successfully"})
}

func GetSymbolAlerts(c* gin.Context){

}
