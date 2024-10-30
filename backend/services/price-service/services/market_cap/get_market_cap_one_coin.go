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

func getMarketCapOneCoin(coinInput *models.CoinInput, context *gin.Context) {
	// Load env variable
	key, err := utils.GetKeyApi()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error server from getting env"})
		return
	}
	// End load env variable

	// Send request to get categories of coinmarketcap server
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		utils.ShowError(http.StatusInternalServerError, "Error from getting new request", context)
		return
	}

	q := url.Values{}
	q.Add("id", fmt.Sprint(coinInput.IdCoin))
	if coinInput.ConvertId != "" {
		// convert to use USD (if want to use BTC change convert_id to 1)
		q.Add("convert_id", "2781")
	}
	q.Add("aux", "num_market_pairs,cmc_rank,date_added,platform,max_supply,circulating_supply,total_supply,is_active,is_fiat")

	req.Header.Add("Accepts", "application/json")
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
