package marketcap

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dath-241/coin-price-be-go/services/price-service/models"
	"github.com/dath-241/coin-price-be-go/services/price-service/utils"
	"github.com/gin-gonic/gin"
)

func getMarketCapItemCategory(categoryItem *models.CategoryInput, context *gin.Context) {
	// Load env variable
	key, err := utils.GetKeyApi()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error server from getting env"})
		return
	}
	// End load env variable

	// Send request to get categories of coinmarketcap server
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/category", nil)
	if err != nil {
		utils.ShowError(http.StatusInternalServerError, "Error from getting new request", context)
		return
	}

	q := url.Values{}
	q.Add("id", categoryItem.IdCategory)
	q.Add("start", fmt.Sprint(categoryItem.Start))
	q.Add("limit", fmt.Sprint(categoryItem.Limit))
	if categoryItem.Convert != "" {
		q.Add("convert", fmt.Sprint(categoryItem.Convert))
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", key)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		utils.ShowError(http.StatusInternalServerError, "Error sending request to server", context)
		return
	}
	defer resp.Body.Close()
	statusCode, _ := strconv.ParseInt(resp.Status, 10, 64)
	respBody, _ := io.ReadAll(resp.Body)
	context.Data(int(statusCode), "application/json", respBody)
}
