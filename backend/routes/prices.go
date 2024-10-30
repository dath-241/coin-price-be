package routes

import (
	fundingrate "github.com/dath-241/coin-price-be-go/services/price-service/services/funding_rate"
	marketcap "github.com/dath-241/coin-price-be-go/services/price-service/services/market_cap"
	"github.com/gin-gonic/gin"
)

func getFundingRate(context *gin.Context) {
	fundingrate.GetFundingRate(context)
}

func getFundingRateCountdown(context *gin.Context) {
	fundingrate.GetFundingRateCountdown(context)
}

func getMarketCap(context *gin.Context) {
	marketcap.GetMarketCap(context)
}
