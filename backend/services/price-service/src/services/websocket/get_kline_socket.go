package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dath-241/coin-price-be-go/services/price-service/src/models"
	"github.com/dath-241/coin-price-be-go/services/price-service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func KlineSocket(context *gin.Context) {

	// Create websocket
	ws, err := Upgrade(context.Writer, context.Request)
	if err != nil {
		log.Println("Upgrade error: ", err)
		return
	}
	defer ws.Close()

	symbol := strings.ToLower(context.Query("symbol"))
	wsURL := fmt.Sprintf("wss://stream.binance.com/stream?streams=%s@kline_1s", symbol)

	headers := http.Header{}
	headers.Add("method", "SUBSCRIBE")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	defer conn.Close()

	done := make(chan struct{})
	// handle symbol error
	isReceivedMessage := make(chan bool)
	timeoutDuration := 5 * time.Second

	go func() {
		defer close(done)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error: ", err)
				return
			}

			// alert get message (symbol not error)
			isReceivedMessage <- true

			var KlineResponse models.KlineWebsocket

			if err = json.Unmarshal(message, &KlineResponse); err != nil {
				log.Println("JSON unmarshal error: ", err)
				continue
			}

			response := processKlineResponse(&KlineResponse)

			responseJSON, err := json.Marshal(&response)
			if err != nil {
				errorMsg := fmt.Sprintf("JSON marshal error: %s", err.Error())
				utils.ShowErrorSocket(ws, errorMsg)
				continue
			}

			if err = ws.WriteMessage(websocket.TextMessage, []byte(responseJSON)); err != nil {
				errorMsg := fmt.Sprintf("Write error to client %s", err.Error())
				utils.ShowErrorSocket(ws, errorMsg)
				return
			}
		}
	}()

	// after 5 seconds, if not response
	go func() {
		for {
			select {
			case <-isReceivedMessage:
				continue
			case <-time.After(timeoutDuration):
				errorMSG := "Symbol error"
				utils.ShowErrorSocket(ws, errorMSG)
				// close connect socket with binance server
				conn.Close()
				// close socket with user
				ws.Close()
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
			log.Println("Disconnecting from Websocket")
			break
		}
	}

	<-done
}

func processKlineResponse(KlineResponse *models.KlineWebsocket) map[string]interface{} {
	return map[string]interface{}{
		"symbol":              KlineResponse.Data.Symbol,
		"eventTime":           utils.ConvertMillisecondsToTimestamp(KlineResponse.Data.EventTime),
		"startTime":           utils.ConvertMillisecondsToTimestamp(KlineResponse.Data.KData.StartTime),
		"closeTime":           utils.ConvertMillisecondsToTimestamp(KlineResponse.Data.KData.CloseTime),
		"openPrice":           KlineResponse.Data.KData.OpenPrice,
		"highPrice":           KlineResponse.Data.KData.HighPrice,
		"lowPrice":            KlineResponse.Data.KData.LowPrice,
		"baseAssetVolume":     KlineResponse.Data.KData.BaseAssetVolume,
		"quoteAssetVolume":    KlineResponse.Data.KData.QuoteAssetVolume,
		"takerBuyBaseVolume":  KlineResponse.Data.KData.TakerBuyBaseVolume,
		"takerBuyQuoteVolume": KlineResponse.Data.KData.TakerBuyQuoteVolume,
	}
}
