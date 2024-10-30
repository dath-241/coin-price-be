package fundingrate

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"
	"github.com/dath-241/coin-price-be-go/services/price-service/utils"
	"github.com/gin-gonic/gin"
)

func GetFundingRateCountdownRealtime(inputFundingRate *models.InputFundingRate, context *gin.Context) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://fapi.binance.com/fapi/v1/premiumIndex", nil)
	if err != nil {
		utils.ShowError(http.StatusInternalServerError, "Error from create new request", context)
		return
	}

	q := url.Values{}
	q.Add("symbol", inputFundingRate.Symbol)

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		utils.ShowError(http.StatusInternalServerError, "Error sending request to server", context)
		return
	}
	defer resp.Body.Close()

	var response models.ResponseFundingRate
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		utils.ShowError(http.StatusInternalServerError, "Error decoding response", context)
		return
	}

	countdown := response.NextFundingTime - response.CurrentTime
	context.JSON(http.StatusOK, gin.H{
		"symbol":          response.Symbol,
		"nextFundingTime": response.NextFundingTime,
		"time":            response.CurrentTime,
		"countdown":       countdown,
	})
}
