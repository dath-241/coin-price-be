package kline

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type KlineResponse struct {
	Symbol    string          `json:"symbol"`
	Interval  string          `json:"interval"`
	EventTime string          `json:"eventTime"`
	KlineData []KLineEachData `json:"kline_data"`
}

type KLineEachData struct {
	Time   string  `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

var mockKlineData = [][]interface{}{
	{
		float64(1689033600000), // "2023-07-10T00:00:00Z"
		"30147.8",              // open
		"31040.0",              // high
		"29928.8",              // low
		"30396.9",              // close
		"429115.537",           // volume
	},
	{
		float64(1689120000000), // "2023-07-11T00:00:00Z"
		"30396.9",              // open
		"30804.9",              // high
		"30261.4",              // low
		"30608.4",              // close
		"298904.747",           // volume
	},
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestGetKline(t *testing.T) {
	router := setupTestRouter()
	router.GET("/kline", GetKline)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
	}{
		{
			name:           "Valid Request",
			queryParams:    "symbol=BTCUSDT&interval=1d",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing Symbol",
			queryParams:    "interval=1d",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Missing Interval",
			queryParams:    "symbol=BTCUSDT",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Empty Parameters",
			queryParams:    "",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/kline?"+tt.queryParams, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetKlineData(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
		interval := r.URL.Query().Get("interval")

		if symbol == "" || interval == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockKlineData)
	}))
	defer mockServer.Close()

	tests := []struct {
		name           string
		symbol         string
		interval       string
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "Successful Request",
			symbol:         "BTCUSDT",
			interval:       "1d",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "Empty Symbol",
			symbol:         "",
			interval:       "1d",
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  false,
		},
		{
			name:           "Empty Interval",
			symbol:         "BTCUSDT",
			interval:       "",
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			GetKlineData(tt.symbol, tt.interval, c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse {
				var response KlineResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check response structure
				assert.Equal(t, tt.symbol, response.Symbol)
				assert.Equal(t, tt.interval, response.Interval)
				assert.NotEmpty(t, response.EventTime)

				// Verify kline data format and content
				assert.NotEmpty(t, response.KlineData)
				if len(response.KlineData) > 0 {
					firstKline := response.KlineData[0]

					// Verify time format
					_, err := time.Parse(time.RFC3339, firstKline.Time)
					assert.NoError(t, err)

					// Verify numeric values
					assert.Equal(t, 30147.8, firstKline.Open)
					assert.Equal(t, 31040.0, firstKline.High)
					assert.Equal(t, 29928.8, firstKline.Low)
					assert.Equal(t, 30396.9, firstKline.Close)
					assert.Equal(t, 429115.537, firstKline.Volume)
				}
			}
		})
	}
}

func TestChangeToFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
	}{
		{
			name:     "Valid Price String",
			input:    "30147.8",
			expected: 30147.8,
		},
		{
			name:     "Valid Volume String",
			input:    "429115.537",
			expected: 429115.537,
		},
		{
			name:     "Invalid String",
			input:    "invalid",
			expected: 0.0,
		},
		{
			name:     "Empty String",
			input:    "",
			expected: 0.0,
		},
		{
			name:     "Non-String Input",
			input:    123.45,
			expected: 0.0,
		},
		{
			name:     "Nil Input",
			input:    nil,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ChangeToFloat(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestKlineResponseStructure(t *testing.T) {
	expectedResponse := KlineResponse{
		Symbol:    "BTCUSDT",
		Interval:  "1d",
		EventTime: "2024-11-20 09:09:02",
		KlineData: []KLineEachData{
			{
				Time:   "2023-07-10T00:00:00Z",
				Open:   30147.8,
				High:   31040.0,
				Low:    29928.8,
				Close:  30396.9,
				Volume: 429115.537,
			},
			{
				Time:   "2023-07-11T00:00:00Z",
				Open:   30396.9,
				High:   30804.9,
				Low:    30261.4,
				Close:  30608.4,
				Volume: 298904.747,
			},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(expectedResponse)
	assert.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledResponse KlineResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	assert.NoError(t, err)

	// Verify structure and values
	assert.Equal(t, expectedResponse.Symbol, unmarshaledResponse.Symbol)
	assert.Equal(t, expectedResponse.Interval, unmarshaledResponse.Interval)
	assert.Equal(t, expectedResponse.EventTime, unmarshaledResponse.EventTime)
	assert.Equal(t, len(expectedResponse.KlineData), len(unmarshaledResponse.KlineData))

	// Verify each kline data entry
	for i, expectedKline := range expectedResponse.KlineData {
		actualKline := unmarshaledResponse.KlineData[i]
		assert.Equal(t, expectedKline.Time, actualKline.Time)
		assert.Equal(t, expectedKline.Open, actualKline.Open)
		assert.Equal(t, expectedKline.High, actualKline.High)
		assert.Equal(t, expectedKline.Low, actualKline.Low)
		assert.Equal(t, expectedKline.Close, actualKline.Close)
		assert.Equal(t, expectedKline.Volume, actualKline.Volume)
	}
}
