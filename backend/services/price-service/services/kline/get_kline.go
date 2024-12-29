package kline

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"
	"github.com/dath-241/coin-price-be-go/services/price-service/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get Kline data
// @Description Fetches Kline data for a specific symbol and interval from Binance API
// @Tags Kline
// @Param Authorization header string true "Authorization token"
// @Param symbol query string true "Symbol for which to fetch Kline data (e.g., BTCUSDT)"
// @Param interval query string true "Interval for Kline data (e.g., 1m, 5m, 1h, 1d)"
// @Success 200 {object} models.ResponseKline "Successful response with Kline data"
// @Failure 400 {object} models.ErrorResponseInputMissing "Missing Data"
// @Failure 500 {object} models.ErrorResponseDataInternalServerError "Internal server error"
// @Router /api/v1/vip1/kline [get]
func GetKline(context *gin.Context) {
	symbol := context.Query("symbol")
	interval := context.Query("interval")
	GetKlineData(symbol, interval, context)
}

func GetKlineData(symbol, interval string, context *gin.Context) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://fapi.binance.com/fapi/v1/klines", nil)
	if err != nil {
		utils.ShowError(http.StatusInternalServerError, err.Error(), context)
		return
	}

	if symbol == "" || interval == "" {
		utils.ShowError(http.StatusBadRequest, "Missing data", context)
		return
	}

	q := url.Values{}
	q.Add("symbol", symbol)
	q.Add("interval", interval)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		utils.ShowError(http.StatusInternalServerError, "Internal server error", context)
		return
	}
	defer resp.Body.Close()

	respStatusCode := resp.StatusCode
	if respStatusCode != http.StatusOK {
		utils.ShowError(http.StatusInternalServerError, "Internal server error", context)
		return
	}

	var data [][]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		utils.ShowError(http.StatusInternalServerError, "Internal server error", context)
		return
	}

	var response models.KlineResponse
	response.UpdateKlineResponse(symbol, interval, utils.GetTimeNow())
	for _, value := range data {
		var kline models.KLineEachData
		timeData := int64(value[0].(float64))
		timeKline := utils.ConvertMilisecondToTimeFormatedRFC3339(timeData)
		open := ChangeToFloat(value[1])
		high := ChangeToFloat(value[2])
		low := ChangeToFloat(value[3])
		close := ChangeToFloat(value[4])
		volume := ChangeToFloat(value[5])
		kline.UpdateKlineEachData(timeKline, open, high, low, close, volume)
		response.UpdateKlineResponseData(&kline)
	}
	// data response
	// [
	// "symbol": "BTCUSDT",
	// "interval": "1d",
	// "eventTime": "2024-11-20 09:09:02",
	// "kline_data": [
	//     {
	//         "time": "2023-07-10T00:00:00Z",
	//         "open": 30147.8,
	//         "high": 31040,
	//         "low": 29928.8,
	//         "close": 30396.9,
	//         "volume": 429115.537
	//     },
	//     {
	//         "time": "2023-07-11T00:00:00Z",
	//         "open": 30396.9,
	//         "high": 30804.9,
	//         "low": 30261.4,
	//         "close": 30608.4,
	//         "volume": 298904.747
	//     },]

	context.JSON(http.StatusOK, response)
}
func ChangeToFloat(data interface{}) float64 {
	if strVal, ok := data.(string); ok {
		result, _ := strconv.ParseFloat(strVal, 64)
		return result
	}
	return 0.0
}
