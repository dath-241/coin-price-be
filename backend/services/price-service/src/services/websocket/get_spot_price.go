package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dath-241/coin-price-be-go/services/price-service/src/models"
	"github.com/dath-241/coin-price-be-go/services/price-service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func SpotPriceSocket(context *gin.Context) {
	ws, err := Upgrade(context.Writer, context.Request)
	if err != nil {
		log.Println("Upgrade error: ", err)
		return
	}
	defer ws.Close()

	symbol := strings.ToLower(context.Query("symbol"))

	wsURL := fmt.Sprintf("wss://stream.binance.com/ws/%s@ticker", symbol)

	headers := http.Header{}
	headers.Add("method", "SUBSCRIBE")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error: ", err)
				return
			}

			var tickerResponse models.SpotTickerWebSocket
			if err = json.Unmarshal(message, &tickerResponse); err != nil {
				log.Println("JSON unmarshal error: ", err)
				continue
			}
			response := map[string]interface{}{
				"symbol":    tickerResponse.Symbol,
				"price":     tickerResponse.LastPrice,
				"eventTime": utils.ConvertMillisecondsToTimestamp(tickerResponse.EventTime),
			}

			responseJSON, err := json.Marshal(&response)
			if err != nil {
				errorMsg := fmt.Sprintf("JSON marshal error: %s", err.Error())
				utils.ShowErrorSocket(ws, errorMsg)
				continue
			}

			if err := ws.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
				errorMsg := fmt.Sprintf("Write error to client %s", err.Error())
				utils.ShowErrorSocket(ws, errorMsg)
				return
			}
		}
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}

		if string(msg) == "disconnect" {
			log.Println("Disconnecting from WebSocket")
			break
		}
	}

	<-done
}
