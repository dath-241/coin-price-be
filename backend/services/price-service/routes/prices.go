package routes

import (
	// fundingrate "github.com/dath-241/coin-price-be-go/services/price-service/services/funding_rate"
	// "github.com/dath-241/coin-price-be-go/services/price-service/services/kline"
	// "github.com/dath-241/coin-price-be-go/services/price-service/services/websocket"
	"github.com/dath-241/coin-price-be-go/services/price-service/services/kline"
	"github.com/dath-241/coin-price-be-go/services/price-service/services/websocket"
	fundingrate "github.com/dath-241/coin-price-be-go/services/price-service/src/services/funding_rate"
	"github.com/gin-gonic/gin"
)

func getFundingRate(context *gin.Context) {
	fundingrate.GetFundingRate(context)
}

func getKline(context *gin.Context) {
	kline.GetKline(context)
}

func getWebsocketFundingRate(context *gin.Context) {
	websocket.FundingRateSocket(context)
}

func getWebsocketKline(context *gin.Context) {
	websocket.KlineSocket(context)
}

func getWebsocketMarketCap(context *gin.Context) {
	websocket.MarketCapSocket(context)
}

func getWebsocketSpotPrice(context *gin.Context) {
	websocket.SpotPriceSocket(context)
}

func getWebsocketFuturePrice(context *gin.Context) {
	websocket.FuturePriceSocket(context)
}
