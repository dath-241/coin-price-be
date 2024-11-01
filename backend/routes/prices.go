package routes

import (
	fundingrate "github.com/dath-241/coin-price-be-go/services/price-service/services/funding_rate"
	"github.com/gin-gonic/gin"
)

func getFundingRate(context *gin.Context) {
	fundingrate.GetFundingRate(context)
}
