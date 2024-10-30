package marketcap

import (
	"net/http"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"
	"github.com/dath-241/coin-price-be-go/services/price-service/utils"
	"github.com/gin-gonic/gin"
)

func GetMarketCap(context *gin.Context) {

	// get query type
	type_query := context.Query("type")
	if type_query == "" {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Your query is missing"})
		return
	}
	// End get query type

	if type_query == "categories" {
		// case 1: list all categories
		// Binding data from body {start, limit}
		var categoriesInp models.CategoriesInput
		err := context.ShouldBindJSON(&categoriesInp)
		if err != nil {
			utils.ShowError(http.StatusInternalServerError, "Error from binding data categories", context)
			return
		}
		getMarketCapCategories(&categoriesInp, context)
		return
	} else if type_query == "category" {
		// case 2: list all coin with categories
		// Binding data from body {id_category, start, limit, convert (USD, BTC)}
		var categoryInp models.CategoryInput
		err := context.ShouldBindJSON(&categoryInp)
		if err != nil {
			utils.ShowError(http.StatusInternalServerError, "Error from binding data category", context)
			return
		}
		// call api to get list item with category id explicit
		getMarketCapItemCategory(&categoryInp, context)
		return
	} else if type_query == "coin" {
		var coinInp models.CoinInput
		err := context.ShouldBindJSON(&coinInp)
		if err != nil {
			utils.ShowError(http.StatusInternalServerError, "Error from binding data coin", context)
			return
		}
		// call api to get market cap of coin with id in category id
		getMarketCapOneCoin(&coinInp, context)
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Not categories"})
}
