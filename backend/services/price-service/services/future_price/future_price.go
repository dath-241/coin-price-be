package future_price

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BinanceFutureResponse struct {
	Symbol               string `json:"symbol"`
	MarkPrice            string `json:"markPrice"`
	IndexPrice           string `json:"indexPrice"`
	EstimatedSettlePrice string `json:"estimatedSettlePrice"`
	LastFundingRate      string `json:"lastFundingRate"`
	NextFundingTime      int64  `json:"nextFundingTime"`
	InterestRate         string `json:"interestRate"`
	Time                 int64  `json:"time"`
}

type FuturePriceResponse struct {
	EventTime string `json:"eventTime"`
	Price     string `json:"price"`
	Symbol    string `json:"symbol"`
}

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
	var binanceResp BinanceFutureResponse
	if err := json.NewDecoder(resp.Body).Decode(&binanceResp); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to decode response: %v", err)})
		return
	}

	// Convert timestamp to formatted date string
	eventTime := time.Unix(binanceResp.Time/1000, 0).Format("2006-01-02 15:04:05")

	// Create our response structure
	response := &FuturePriceResponse{
		EventTime: eventTime,
		Price:     binanceResp.MarkPrice,
		Symbol:    binanceResp.Symbol,
	}

	ctx.JSON(http.StatusOK, response)
}
