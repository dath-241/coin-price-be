package spot_price

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"

	"github.com/gin-gonic/gin"
)

// @Summary Get real-time spot price data
// @Description Retrieves current spot price information for a specified trading pair from Binance Spot
// @Tags Spot price
// @Produce json
// @Param symbol query string true "Trading pair symbol (e.g., BTCUSDT)" example("BTCUSDT")
// @Success 200 {object} models.ResponseSpotPrice "Successful response with spot price data"
// @Failure 400 {object} models.ErrorResponseDataMissing "Invalid symbol or request parameters"
// @Failure 404 {object} models.ErrorResponseDataNotFound "Symbol not found"
// @Failure 500 {object} models.ErrorResponseDataInternalServerError "Failed to fetch price"
// @Router /api/v1/spot-price [get]
func GetSpotPrice(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol cannot be empty"})
		return
	}

	// Construct the Binance API URL
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v2/ticker/price?symbol=%s", symbol)

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch price: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Check if the response status is not OK
	if resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("API returned status code: %d", resp.StatusCode)})
		return
	}

	// Decode the Binance response
	var binanceResp models.ResponseBinance
	if err := json.NewDecoder(resp.Body).Decode(&binanceResp); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to decode response: %v", err)})
		return
	}

	// Convert timestamp to formatted date string
	eventTime := time.Unix(binanceResp.Time/1000, 0).Format("2006-01-02 15:04:05")

	// Create our response structure
	response := &models.ResponseSpotPrice{
		EventTime: eventTime,
		Price:     binanceResp.Price,
		Symbol:    binanceResp.Symbol,
	}

	ctx.JSON(http.StatusOK, response)
}
