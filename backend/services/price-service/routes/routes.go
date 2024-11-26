package routes

import (
	docs "github.com/dath-241/coin-price-be-go/docs"
	middlewares "github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
	"github.com/dath-241/coin-price-be-go/services/price-service/services/future_price"
	"github.com/dath-241/coin-price-be-go/services/price-service/services/spot_price"
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
	server.GET("/swagger/price/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	authenticated := server.Group("/api")
	// Funding rate
	authenticated.GET("/v1/funding-rate", getFundingRate)
	authenticated.GET("/v1/funding-rate/websocket", getWebsocketFundingRate)
	// Spot price
	authenticated.GET("/v1/spot-price", spot_price.GetSpotPrice)
	authenticated.GET("/v1/spot-price/websocket", getWebsocketSpotPrice)
	// Future price
	authenticated.GET("/v1/future-price", future_price.GetFuturePrice)
	authenticated.GET("/v1/future-price/websocket", getWebsocketFuturePrice)
	// Market stats
	authenticated.GET("/v1/market-stats", getWebsocketMarketCap)
	// Kline
	authenticated.GET("/v1/vip1/kline", middlewares.AuthMiddleware("VIP-1", "VIP-2", "VIP-3"), getKline)
	authenticated.GET("/v1/vip1/kline/websocket", middlewares.AuthMiddleware("VIP-1", "VIP-2", "VIP-3"), getWebsocketKline)
}
