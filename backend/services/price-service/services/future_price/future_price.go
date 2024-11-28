package future_price

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"
	"github.com/gin-gonic/gin"
)

// @Summary Get real-time future price data
// @Description Retrieves current future price information for a specified trading pair from Binance Futures
// @Tags Future price
// @Produce json
// @Param symbol query string true "Trading pair symbol (e.g., BTCUSDT)" example("BTCUSDT")
// @Success 200 {object} models.ResponseFuturePrice "Successful response with future price data"
// @Failure 400 {object} models.ErrorResponseDataMissing "Invalid symbol or request parameters"
// @Failure 404 {object} models.ErrorResponseDataNotFound "Symbol not found"
// @Failure 500 {object} models.ErrorResponseDataInternalServerError "Failed to fetch price"
// @Router /api/v1/future-price [get]
func GetFuturePrice(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol cannot be empty"})
		return
	}

	// Construct the Binance Futures API URL
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/premiumIndex?symbol=%s", symbol)

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch future price: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Check if the response status is not OK
	if resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("API returned status code: %d", resp.StatusCode)})
		return
	}

	// Decode the Binance response
	var binanceResp models.ResponseBinanceFuture
	if err := json.NewDecoder(resp.Body).Decode(&binanceResp); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to decode response: %v", err)})
		return
	}

	// Convert timestamp to formatted date string
	eventTime := time.Unix(binanceResp.Time/1000, 0).Format("2006-01-02 15:04:05")

	// Create our response structure
	response := &models.ResponseFuturePrice{
		EventTime: eventTime,
		Price:     binanceResp.MarkPrice,
		Symbol:    binanceResp.Symbol,
	}

	ctx.JSON(http.StatusOK, response)
}
