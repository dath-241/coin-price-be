package websocket

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestMarketCapSocket(t *testing.T) {
	// Setup Gin router
	router := gin.Default()
	router.GET("/ws/market-cap", MarketCapSocket)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Convert http://... to ws://...
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/market-cap?symbol=bitcoin"

	// Connect to websocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Read the first message
	_, msg, err := ws.ReadMessage()
	assert.NoError(t, err)

	// Parse response
	var response models.FormatMarketCapResponse
	err = json.Unmarshal(msg, &response)
	assert.NoError(t, err)

	// Assert response structure
	assert.Equal(t, "btc", response.Symbol)
	assert.Greater(t, response.MarketCap, int64(0))
	assert.Greater(t, response.TotalVolume, int64(0))

	// Test disconnect
	err = ws.WriteMessage(websocket.TextMessage, []byte("disconnect"))
	assert.NoError(t, err)
}

func TestCreateResponseFormat(t *testing.T) {
	// Test data
	symbol := "btc"
	marketCap := int64(1000000000)
	totalVolume := int64(500000000)

	// Create response
	response := models.CreateReponseFormat(symbol, marketCap, totalVolume)

	// Assert values
	assert.Equal(t, symbol, response.Symbol)
	assert.Equal(t, marketCap, response.MarketCap)
	assert.Equal(t, totalVolume, response.TotalVolume)
}
