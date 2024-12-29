package fundingrate

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockFundingRateFirst represents mock data for first API response
var MockFundingRateFirst = models.FundingRateFirst{
	Symbol:          "BTCUSDT",
	FundingRate:     "0.0001",
	NextFundingTime: time.Now().Add(time.Hour*8).UnixNano() / int64(time.Millisecond),
	EventTime:       time.Now().UnixNano() / int64(time.Millisecond),
}

// MockFundingRateSecond represents mock data for second API response
var MockFundingRateSecond = models.FundingRateSecond{
	Symbol:                   "BTCUSDT",
	AdjustedFundingRateCap:   "0.0075",
	AdjustedFundingRateFloor: "-0.0075",
	FundingIntervalHours:     8,
}

// Expected response after processing
var ExpectedResponse = models.ResponseFundingRate{
	Symbol:                   "BTCUSDT",
	FundingRate:              "0.0001",
	FundingCountDown:         "07:59:59", // Example countdown
	EventTime:                time.Now().Format("2006-01-02 15:04:05"),
	AdjustedFundingRateCap:   "0.0075",
	AdjustedFundingRateFloor: "-0.0075",
	FundingIntervalHours:     8,
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestGetFundingRate(t *testing.T) {
	router := setupTestRouter()
	router.GET("/funding-rate", GetFundingRate)

	tests := []struct {
		name           string
		symbol         string
		expectedStatus int
	}{
		{
			name:           "Valid symbol",
			symbol:         "BTCUSDT",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing symbol",
			symbol:         "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/funding-rate?symbol="+tt.symbol, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response models.ResponseFundingRate
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.symbol, response.Symbol)
				assert.NotEmpty(t, response.FundingRate)
				assert.NotEmpty(t, response.FundingCountDown)
				assert.NotEmpty(t, response.EventTime)
			}
		})
	}
}

func setupMockServer(t *testing.T, firstEndpoint bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		if firstEndpoint {
			assert.Contains(t, r.URL.Path, "/fapi/v1/premiumIndex")
			json.NewEncoder(w).Encode(MockFundingRateFirst)
		} else {
			assert.Contains(t, r.URL.Path, "/fapi/v1/fundingInfo")
			json.NewEncoder(w).Encode([]models.FundingRateSecond{MockFundingRateSecond})
		}
	}))
}

func TestGetDataFundingFirst(t *testing.T) {
	server := setupMockServer(t, true)
	defer server.Close()

	tests := []struct {
		name         string
		symbol       string
		expectedCode models.StatusCode
		expectError  bool
	}{
		{
			name:         "Valid symbol",
			symbol:       "BTCUSDT",
			expectedCode: models.StatusCode(200),
			expectError:  false,
		},
		{
			name:         "Invalid symbol",
			symbol:       "INVALID",
			expectedCode: models.StatusCode(400),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, statusCode, err := GetDataFundingFirst(tt.symbol)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, statusCode)
				assert.NotNil(t, response)
				assert.Equal(t, tt.symbol, response.Symbol)
				assert.NotEmpty(t, response.FundingRate)
				assert.Greater(t, response.NextFundingTime, response.EventTime)
			}
		})
	}
}

func TestGetDataFundingSecond(t *testing.T) {
	server := setupMockServer(t, false)
	defer server.Close()

	tests := []struct {
		name         string
		symbol       string
		expectedCode models.StatusCode
	}{
		{
			name:         "Valid symbol",
			symbol:       "QTUMUSDT",
			expectedCode: models.StatusCode(200),
		},
		{
			name:         "Non-existent symbol",
			symbol:       "NONEXISTENT",
			expectedCode: models.StatusCode(404),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, statusCode := GetDataFundingSecond(tt.symbol)

			assert.Equal(t, tt.expectedCode, statusCode)
			if statusCode == http.StatusOK {
				assert.NotNil(t, response)
				assert.Equal(t, tt.symbol, response.Symbol)
				assert.NotEmpty(t, response.AdjustedFundingRateCap)
				assert.NotEmpty(t, response.AdjustedFundingRateFloor)
				assert.Greater(t, response.FundingIntervalHours, 0)
			} else {
				assert.Nil(t, response)
			}
		})
	}
}

func TestProcessResponse(t *testing.T) {
	currentTime := time.Now()
	nextFundingTime := currentTime.Add(time.Hour * 8)

	resp1 := &models.FundingRateFirst{
		Symbol:          "BTCUSDT",
		FundingRate:     "0.0001",
		NextFundingTime: nextFundingTime.UnixNano() / int64(time.Millisecond),
		EventTime:       currentTime.UnixNano() / int64(time.Millisecond),
	}

	resp2 := &models.FundingRateSecond{
		Symbol:                   "BTCUSDT",
		AdjustedFundingRateCap:   "0.0075",
		AdjustedFundingRateFloor: "-0.0075",
		FundingIntervalHours:     8,
	}

	var result models.ResponseFundingRate

	t.Run("Process valid responses", func(t *testing.T) {
		ProcessResponse(resp1, resp2, &result)

		// Verify all fields match the ResponseFundingRate struct
		assert.Equal(t, resp1.Symbol, result.Symbol)
		assert.Equal(t, resp1.FundingRate, result.FundingRate)
		assert.Equal(t, resp2.AdjustedFundingRateCap, result.AdjustedFundingRateCap)
		assert.Equal(t, resp2.AdjustedFundingRateFloor, result.AdjustedFundingRateFloor)
		assert.Equal(t, resp2.FundingIntervalHours, result.FundingIntervalHours)

		// Verify countdown format (HH:MM:SS)
		assert.Regexp(t, `^\d{2}:\d{2}:\d{2}$`, result.FundingCountDown)

		// Verify event time format (YYYY-MM-DD HH:MM:SS)
		assert.Regexp(t, `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, result.EventTime)
	})

	t.Run("Process with unknown funding info", func(t *testing.T) {
		resp2Unknown := &models.FundingRateSecond{
			Symbol:                   "BTCUSDT",
			AdjustedFundingRateCap:   "unknown",
			AdjustedFundingRateFloor: "unknown",
			FundingIntervalHours:     -1,
		}

		var resultUnknown models.ResponseFundingRate
		ProcessResponse(resp1, resp2Unknown, &resultUnknown)

		assert.Equal(t, "unknown", resultUnknown.AdjustedFundingRateCap)
		assert.Equal(t, "unknown", resultUnknown.AdjustedFundingRateFloor)
		assert.Equal(t, -1, resultUnknown.FundingIntervalHours)
	})
}
