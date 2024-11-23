package spot_price

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BinanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int64  `json:"time"`
}

type SpotPriceResponse struct {
	EventTime string `json:"eventTime"`
	Price     string `json:"price"`
	Symbol    string `json:"symbol"`
}

// @Summary Get real-time spot price data
// @Description Retrieves current spot price information for a specified trading pair from Binance Futures
// @Tags Spot price
// @Produce json
// @Param symbol query string true "Trading pair symbol (e.g., BTCUSDT)" example("BTCUSDT")
// @Success 200 {object} SpotPriceResponse "Successful response with funding rate data"
// @Failure 400 {object} SpotPriceResponse "Invalid symbol or request parameters"
// @Failure 404 {object} SpotPriceResponse "Symbol not found"
// @Failure 500 {object} SpotPriceResponse "failed to fetch price"
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
	var binanceResp BinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&binanceResp); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to decode response: %v", err)})
		return
	}

	// Convert timestamp to formatted date string
	eventTime := time.Unix(binanceResp.Time/1000, 0).Format("2006-01-02 15:04:05")

	// Create our response structure
	response := &SpotPriceResponse{
		EventTime: eventTime,
		Price:     binanceResp.Price,
		Symbol:    binanceResp.Symbol,
	}

	ctx.JSON(http.StatusOK, response)
}
