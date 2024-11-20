package routes

import (
	docs "github.com/dath-241/coin-price-be-go/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Your API Title
// @version         1.0
// @description     Your API Description
// @host            localhost:8080
// @BasePath        /api
func RegisterRoutes(server *gin.Engine) {

	docs.SwaggerInfo.BasePath = "/api"

	// Move swagger route outside of the authenticated group
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	authenticated := server.Group("/api")
	authenticated.GET("/v1/funding-rate", getFundingRate)
	authenticated.GET("/v1/vip1/kline", getKline)
	authenticated.GET("/v1/funding-rate/websocket", getWebsocketFundingRate)
	authenticated.GET("/v1/vip1/kline/websocket", getWebsocketKline)
	authenticated.GET("/v1/spot-price/websocket", getWebsocketSpotPrice)
	authenticated.GET("/v1/future-price/websocket", getWebsocketFuturePrice)
	authenticated.GET("/v1/market-stats", getWebsocketMarketCap)
}
