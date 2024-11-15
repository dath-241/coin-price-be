package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	models "github.com/dath-241/coin-price-be-go/services/trigger-service/models/alert"
)

// fetchSymbolsFromBinance fetches symbols from Binance's API
func FetchSymbolsFromBinance() ([]string, []string, error) {
	url := "https://api.binance.com/api/v3/exchangeInfo"
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("binance API returned status %d", resp.StatusCode)
	}

	var data struct {
		Symbols []models.Symbol `json:"symbols"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil, err
	}

	newSymbols := []string{}
	for _, s := range data.Symbols {
		if s.Status == "TRADING" {
			newSymbols = append(newSymbols, s.Symbol)
		}
	}

	delistedSymbols := []string{}
	for _, s := range data.Symbols {
		if s.Status != "TRADING" {
			delistedSymbols = append(delistedSymbols, s.Symbol)
		}
	}

	return newSymbols, delistedSymbols, nil
}

type PriceResponse struct {
	Price string `json:"price"`
}

type FundingRateResponse struct {
	FundingRate string `json:"fundingRate"`
}

type FuturePriceResponse struct {
	LastPrice string `json:"lastPrice"`
}

// Hàm lấy giá Spot
func GetSpotPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %v", err)
	}

	log.Println("Spot price", symbol, ":", price)
	return price, nil
}

// Hàm lấy Funding Rate
func GetFundingRate(symbol string) (float64, error) {
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/fundingRate?symbol=%s&limit=1", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var results []FundingRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, err
	}

	if len(results) == 0 {
		return 0, fmt.Errorf("no funding rate data available")
	}

	fundingRate, err := strconv.ParseFloat(results[0].FundingRate, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse funding rate: %v", err)
	}

	log.Println("Funding rate", symbol, ":", fundingRate)
	return fundingRate, nil
}

// Hàm lấy giá Future
func GetFuturePrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/ticker/24hr?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result FuturePriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(result.LastPrice, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %v", err)
	}

	log.Println("Future price", symbol, ":", price)
	return price, nil
}
