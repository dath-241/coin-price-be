package spot_price

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetSpotPrice(t *testing.T) {
	// Switch to test mode so we don't get debug output
	gin.SetMode(gin.TestMode)

	// Save the original http.Client
	originalClient := http.DefaultClient

	tests := []struct {
		name           string
		symbol         string
		mockResponse   *models.ResponseBinance
		mockStatusCode int
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful price fetch",
			symbol: "BTCUSDT",
			mockResponse: &models.ResponseBinance{
				Symbol: "BTCUSDT",
				Price:  "50000.00",
				Time:   1677721200000, // 2023-03-02 12:00:00
			},
			mockStatusCode: http.StatusOK,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty symbol",
			symbol:         "",
			mockStatusCode: http.StatusOK,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "symbol cannot be empty",
		},
		{
			name:           "binance api error",
			symbol:         "BTCUSDT",
			mockStatusCode: http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "API returned status code: 500",
		},
		{
			name:           "invalid response body",
			symbol:         "BTCUSDT",
			mockStatusCode: http.StatusOK,
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "failed to decode response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create a mock HTTP server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != nil {
					json.NewEncoder(w).Encode(tt.mockResponse)
				} else {
					w.Write([]byte("invalid json"))
				}
			}))
			defer mockServer.Close()

			// Create a custom http.Client that redirects requests to our mock server
			http.DefaultClient = &http.Client{
				Transport: &mockTransport{
					mockServer: mockServer,
				},
			}

			// Set query parameter
			if tt.symbol != "" {
				c.Request, _ = http.NewRequest(http.MethodGet, "/?symbol="+tt.symbol, nil)
			} else {
				c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
			}

			// Call the function
			GetSpotPrice(c)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				// Assert error message
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				// Assert successful response
				assert.NotNil(t, response["eventTime"])
				assert.Equal(t, tt.mockResponse.Price, response["price"])
				assert.Equal(t, tt.mockResponse.Symbol, response["symbol"])
			}
		})
	}

	// Restore the original http.Client
	http.DefaultClient = originalClient
}

// mockTransport is a custom http.RoundTripper that redirects all requests to our mock server
type mockTransport struct {
	mockServer *httptest.Server
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Replace the request URL with our mock server URL, maintaining the path and query
	mockURL := t.mockServer.URL + req.URL.Path + "?" + req.URL.RawQuery
	newReq, err := http.NewRequest(req.Method, mockURL, req.Body)
	if err != nil {
		return nil, err
	}

	// Copy headers
	newReq.Header = req.Header

	// Use the mock server's client to do the actual request
	return t.mockServer.Client().Do(newReq)
}
