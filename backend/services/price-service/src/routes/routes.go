package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/api")
	authenticated.GET("/v1/funding-rate", getFundingRate)
	authenticated.GET("/v1/vip1/kline", getKline)
	authenticated.GET("/v1/funding-rate/websocket", getWebsocketFundingRate)
	authenticated.GET("/v1/vip1/kline/websocket", getWebsocketKline)
	authenticated.GET("/v1/spot-price/websocket", getWebsocketSpotPrice)
	authenticated.GET("/v1/future-price/websocket", getWebsocketFuturePrice)
	authenticated.GET("/v1/market-stats", getWebsocketMarketCap)
}
