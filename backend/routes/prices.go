package routes

import (
	fundingrate "github.com/dath-241/coin-price-be-go/services/price-service/services/funding_rate"
	"github.com/dath-241/coin-price-be-go/services/price-service/services/kline"
	"github.com/dath-241/coin-price-be-go/services/price-service/services/websocket"
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
