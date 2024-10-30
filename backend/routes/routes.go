package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/api")
	authenticated.GET("/get-funding-rate", getFundingRate)
	authenticated.GET("/get-funding-rate-countdown", getFundingRateCountdown)
	authenticated.GET("/get-marketcap", getMarketCap)
}
