package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/dath-241/coin-price-be-go/services/price-service/models/alert"
)

// fetchSymbolsFromBinance fetches symbols from Binance's API
func fetchSymbolsFromBinance() ([]string, []string, error) {
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
