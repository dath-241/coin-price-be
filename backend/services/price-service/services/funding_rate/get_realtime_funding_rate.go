package fundingrate

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"
	"github.com/dath-241/coin-price-be-go/services/price-service/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get real-time funding rate data
// @Description Retrieves current funding rate information for a specified trading pair from Binance Futures
// @Tags Funding Rate
// @Accept json
// @Produce json
// @Param symbol query string true "Trading pair symbol (e.g., QTUMUSDT)" example("QTUMUSDT")
// @Success 200 {string}  "Successful response with funding rate data"
// @Failure 400 {string}  "Invalid symbol or request parameters"
// @Failure 404 {string}  "Symbol not found"
// @Failure 500 {string}  "Internal server error"
// @Router /v1/funding-rate [get]
func GetFundingRateRealTime(symbol string, context *gin.Context) {
	var responseApi models.ResponseFundingRate
	// get symbol, funding rate, eventTime, countdown
	response1, statusCode, err := GetDataFundingFirst(symbol)
	if err != nil {
		utils.ShowError(int64(statusCode), err.Error(), context)
		return
	}
	// get adjustedFundingRateCap, adjustedFundingRateFloor, fundingInterval if exist
	response2, statusCode := GetDataFundingSecond(symbol)
	if statusCode != http.StatusOK {
		response2 = &models.FundingRateSecond{
			Symbol:                   symbol,
			AdjustedFundingRateCap:   "unknown",
			AdjustedFundingRateFloor: "unknown",
			FundingIntervalHours:     -1,
		}
	}
	ProcessResponse(response1, response2, &responseApi)
	context.JSON(http.StatusOK, responseApi)
}

func GetDataFundingFirst(symbol string) (*models.FundingRateFirst, models.StatusCode, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://fapi.binance.com/fapi/v1/premiumIndex", nil)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("error from request api.")
	}

	q := url.Values{}
	q.Add("symbol", symbol)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("error sending request to server.")
	}
	defer resp.Body.Close()

	respStatusCode := resp.StatusCode
	if respStatusCode == http.StatusBadRequest {
		return nil, http.StatusBadRequest, errors.New("Error information.")
	} else if respStatusCode != http.StatusOK {
		return nil, http.StatusInternalServerError, errors.New("Server error.")
	}

	var response models.FundingRateFirst
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("Error decoding response.")
	}
	return &response, http.StatusOK, nil
}

func GetDataFundingSecond(symbol string) (*models.FundingRateSecond, models.StatusCode) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://fapi.binance.com/fapi/v1/fundingInfo", nil)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	respStatusCode := resp.StatusCode
	if respStatusCode == http.StatusBadRequest {
		return nil, http.StatusBadRequest
	} else if respStatusCode != http.StatusOK {
		return nil, http.StatusInternalServerError
	}

	var response []models.FundingRateSecond
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	// find trading exist
	position := -1
	for index, value := range response {
		if value.Symbol == symbol {
			position = index
			break
		}
	}
	// if not exist, return nil
	if position == -1 {
		return nil, http.StatusNotFound
	}
	// exist return this trading
	return &response[position], http.StatusOK
}

func ProcessResponse(resp1 *models.FundingRateFirst, resp2 *models.FundingRateSecond, result *models.ResponseFundingRate) {

	symbol := resp1.Symbol
	fRate := resp1.FundingRate
	fCountDown := utils.ConvertMillisecondsToHHMMSS(resp1.NextFundingTime - resp1.EventTime)
	time := utils.ConvertMillisecondsToTimestamp(resp1.EventTime)
	fRateCap := resp2.AdjustedFundingRateCap
	fRateFloor := resp2.AdjustedFundingRateFloor
	fIntervalHours := resp2.FundingIntervalHours

	result.UpdateData(symbol, fRate, fCountDown, time, fRateCap, fRateFloor, fIntervalHours)
}
