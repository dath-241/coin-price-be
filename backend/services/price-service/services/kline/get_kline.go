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

	q := url.Values{}
	q.Add("symbol", symbol)
	q.Add("interval", interval)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		utils.ShowError(http.StatusInternalServerError, err.Error(), context)
		return
	}
	defer resp.Body.Close()

	respStatusCode := resp.StatusCode
	if respStatusCode != http.StatusOK {
		context.JSON(respStatusCode, gin.H{"message": "Error get data"})
		return
	}

	var data [][]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error binding data"})
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

	context.JSON(http.StatusOK, response)
}
func ChangeToFloat(data interface{}) float64 {
	if strVal, ok := data.(string); ok {
		result, _ := strconv.ParseFloat(strVal, 64)
		return result
	}
	return 0.0
}
