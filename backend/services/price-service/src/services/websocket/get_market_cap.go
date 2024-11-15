package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dath-241/coin-price-be-go/services/price-service/src/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func MarketCapSocket(context *gin.Context) {
	// Create websocket
	ws, err := Upgrade(context.Writer, context.Request)
	if err != nil {
		log.Println("Upgrade error: ", err)
		return
	}
	defer ws.Close()

	symbol := strings.ToLower(context.Query("symbol"))
	urlMarketCap := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s", symbol)

	// done chan to check if the main go routine is continue or not
	done := make(chan struct{})
	// exit chan to check if the go func is continue or not (check for stop loop)
	exit := make(chan struct{})

	go func() {
		defer close(done)
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		isContinue := processMarketCapSocket(urlMarketCap, ws)
		if !isContinue {
			return
		}
		for {
			select {
			case <-exit:
				return
			case <-ticker.C:
				isContinue := processMarketCapSocket(urlMarketCap, ws)
				if isContinue {
					continue
				} else {
					return
				}
			case <-done:
				return
			}
		}
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			close(exit)
			break
		}
		if string(msg) == "disconnect" {
			close(exit)
			break
		}
	}

	<-done
}

func processMarketCapSocket(urlMarketCap string, ws *websocket.Conn) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlMarketCap, nil)
	if err != nil {
		log.Println("ERROR Request")
		return true
	}

	q := url.Values{}
	q.Add("localization", "false")
	q.Add("tickers", "false")
	q.Add("community_data", "false")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error http request")
		return true
	}

	// Get status code
	statusCode := resp.StatusCode
	if statusCode == http.StatusTooManyRequests {
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Rate limit, please wait."))
		return false

	} else if statusCode != http.StatusOK {
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Symbol missing or invalid"))
		return false
	}

	// get data response
	var marketCapResponse models.MarketCapResponse
	err = json.NewDecoder(resp.Body).Decode(&marketCapResponse)
	resp.Body.Close()
	if err != nil {
		log.Println("Error from getting response, error: ", err)
		return true
	}

	// format response
	dataResponse := models.CreateReponseFormat(
		marketCapResponse.Symbol,
		marketCapResponse.MarketData.MarketCap.USD,
		marketCapResponse.MarketData.TotalVolume.USD,
	)

	responseJSON, err := json.Marshal(&dataResponse)
	if err != nil {
		errorMsg := fmt.Sprintf("JSON marshal error: %s", err.Error())
		log.Println(errorMsg)
		return true
	}

	// return message response to client
	if err = ws.WriteMessage(websocket.TextMessage, []byte(responseJSON)); err != nil {
		errorMsg := fmt.Sprintf("Write error to client %s", err.Error())
		log.Println(errorMsg)
		return false
	}
	return true
}
