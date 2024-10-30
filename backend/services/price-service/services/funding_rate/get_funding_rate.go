package fundingrate

import (
	"net/http"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"
	"github.com/dath-241/coin-price-be-go/services/price-service/utils"
	"github.com/gin-gonic/gin"
)

func GetFundingRate(context *gin.Context) {
	var inputFundingRate models.InputFundingRate
	err := context.ShouldBindJSON(&inputFundingRate)

	if err != nil {
		utils.ShowError(http.StatusInternalServerError, "Error from binding data query with funding rate", context)
		return
	}
	GetRealtimeFundingRate(&inputFundingRate, context)
}

func GetFundingRateCountdown(context *gin.Context) {
	var inputFundingRate models.InputFundingRate
	err := context.ShouldBindJSON(&inputFundingRate)

	if err != nil {
		utils.ShowError(http.StatusInternalServerError, "Error from binding data query with funding rate", context)
		return
	}
	GetFundingRateCountdownRealtime(&inputFundingRate, context)
}
