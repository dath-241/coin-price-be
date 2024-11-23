package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type MockWebsocketServer struct {
	*httptest.Server
	URL string
}

func newMockWebsocketServer() *MockWebsocketServer {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		// Mock Binance server behavior
		klineData := models.KlineWebsocket{
			Data: struct {
				EventType string `json:"e"`
				EventTime int64  `json:"E"`
				Symbol    string `json:"s"`
				KData     struct {
					StartTime           int64  `json:"t"`
					CloseTime           int64  `json:"T"`
					LastTrade           int64  `json:"L"`
					OpenPrice           string `json:"o"`
					ClosePrice          string `json:"c"`
					HighPrice           string `json:"h"`
					LowPrice            string `json:"l"`
					BaseAssetVolume     string `json:"v"`
					QuoteAssetVolume    string `json:"q"`
					TakerBuyBaseVolume  string `json:"V"`
					TakerBuyQuoteVolume string `json:"Q"`
				} `json:"k"`
			}{
				EventType: "kline",
				EventTime: time.Now().UnixMilli(),
				Symbol:    "BTCUSDT",
				KData: struct {
					StartTime           int64  `json:"t"`
					CloseTime           int64  `json:"T"`
					LastTrade           int64  `json:"L"`
					OpenPrice           string `json:"o"`
					ClosePrice          string `json:"c"`
					HighPrice           string `json:"h"`
					LowPrice            string `json:"l"`
					BaseAssetVolume     string `json:"v"`
					QuoteAssetVolume    string `json:"q"`
					TakerBuyBaseVolume  string `json:"V"`
					TakerBuyQuoteVolume string `json:"Q"`
				}{
					StartTime:           time.Now().UnixMilli(),
					CloseTime:           time.Now().Add(time.Second).UnixMilli(),
					OpenPrice:           "50000.00",
					ClosePrice:          "50100.00",
					HighPrice:           "50200.00",
					LowPrice:            "49900.00",
					BaseAssetVolume:     "10.5",
					QuoteAssetVolume:    "525000.00",
					TakerBuyBaseVolume:  "6.3",
					TakerBuyQuoteVolume: "315000.00",
				},
			},
		}

		messageBytes, _ := json.Marshal(klineData)
		conn.WriteMessage(websocket.TextMessage, messageBytes)

		// Read client messages
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}))

	return &MockWebsocketServer{
		Server: mockServer,
		URL:    "ws" + strings.TrimPrefix(mockServer.URL, "http"),
	}
}

func TestKlineSocket(t *testing.T) {
	// Setup mock server
	mockServer := newMockWebsocketServer()
	defer mockServer.Close()

	// Setup Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock HTTP request with query parameters
	req := httptest.NewRequest("GET", "/?symbol=btcusdt", nil)
	c.Request = req

	// Create a done channel to signal test completion
	done := make(chan bool)

	// Run the test
	t.Run("Test successful websocket connection", func(t *testing.T) {
		go func() {
			KlineSocket(c)
			done <- true
		}()

		// Wait for either completion or timeout
		select {
		case <-done:
			// Test passed
		case <-time.After(6 * time.Second):
			t.Error("Test timed out")
		}
	})
}

func TestKlineSocketWithInvalidSymbol(t *testing.T) {
	// Setup mock server
	mockServer := newMockWebsocketServer()
	defer mockServer.Close()

	// Setup Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock HTTP request with invalid symbol
	req := httptest.NewRequest("GET", "/?symbol=invalid", nil)
	c.Request = req

	// Create a done channel to signal test completion
	done := make(chan bool)

	// Run the test
	t.Run("Test invalid symbol connection", func(t *testing.T) {
		go func() {
			KlineSocket(c)
			done <- true
		}()

		// Wait for either completion or timeout
		select {
		case <-done:
			// Test passed
		case <-time.After(6 * time.Second):
			t.Error("Test timed out")
		}
	})
}

func TestProcessKlineResponse(t *testing.T) {
	t.Run("Process valid kline response", func(t *testing.T) {
		input := &models.KlineWebsocket{
			Data: struct {
				EventType string `json:"e"`
				EventTime int64  `json:"E"`
				Symbol    string `json:"s"`
				KData     struct {
					StartTime           int64  `json:"t"`
					CloseTime           int64  `json:"T"`
					LastTrade           int64  `json:"L"`
					OpenPrice           string `json:"o"`
					ClosePrice          string `json:"c"`
					HighPrice           string `json:"h"`
					LowPrice            string `json:"l"`
					BaseAssetVolume     string `json:"v"`
					QuoteAssetVolume    string `json:"q"`
					TakerBuyBaseVolume  string `json:"V"`
					TakerBuyQuoteVolume string `json:"Q"`
				} `json:"k"`
			}{
				Symbol:    "BTCUSDT",
				EventTime: time.Now().UnixMilli(),
				KData: struct {
					StartTime           int64  `json:"t"`
					CloseTime           int64  `json:"T"`
					LastTrade           int64  `json:"L"`
					OpenPrice           string `json:"o"`
					ClosePrice          string `json:"c"`
					HighPrice           string `json:"h"`
					LowPrice            string `json:"l"`
					BaseAssetVolume     string `json:"v"`
					QuoteAssetVolume    string `json:"q"`
					TakerBuyBaseVolume  string `json:"V"`
					TakerBuyQuoteVolume string `json:"Q"`
				}{
					StartTime:           time.Now().UnixMilli(),
					CloseTime:           time.Now().Add(time.Second).UnixMilli(),
					OpenPrice:           "50000.00",
					HighPrice:           "50200.00",
					LowPrice:            "49900.00",
					BaseAssetVolume:     "10.5",
					QuoteAssetVolume:    "525000.00",
					TakerBuyBaseVolume:  "6.3",
					TakerBuyQuoteVolume: "315000.00",
				},
			},
		}

		result := processKlineResponse(input)

		assert.Equal(t, "BTCUSDT", result["symbol"])
		assert.Equal(t, "50000.00", result["openPrice"])
		assert.Equal(t, "50200.00", result["highPrice"])
		assert.Equal(t, "49900.00", result["lowPrice"])
		assert.Equal(t, "10.5", result["baseAssetVolume"])
		assert.Equal(t, "525000.00", result["quoteAssetVolume"])
		assert.Equal(t, "6.3", result["takerBuyBaseVolume"])
		assert.Equal(t, "315000.00", result["takerBuyQuoteVolume"])
	})
}
