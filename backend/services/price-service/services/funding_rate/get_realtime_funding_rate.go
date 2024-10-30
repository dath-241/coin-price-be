package fundingrate

import (
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"
	"github.com/dath-241/coin-price-be-go/services/price-service/utils"
	"github.com/gin-gonic/gin"
)

func GetRealtimeFundingRate(inputFundingRate *models.InputFundingRate, context *gin.Context) {
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

	statusCode, _ := strconv.ParseInt(resp.Status, 10, 64)
	respBody, _ := io.ReadAll(resp.Body)
	context.Data(int(statusCode), "application/json", respBody)
}
