package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	lastSymbols = make(map[string]bool)
	mutex       sync.Mutex
)

// fetchSymbolsFromBinance fetches symbols from Binance's API
func fetchSymbolsFromBinance() (map[string]bool, error) {
	url := "https://api.binance.com/api/v3/exchangeInfo"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance API returned status %d", resp.StatusCode)
	}

	var data struct {
		Symbols []struct {
			Symbol string `json:"symbol"`
			Status string `json:"status"`
		} `json:"symbols"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	symbols := make(map[string]bool)
	for _, s := range data.Symbols {
		if s.Status == "TRADING" {
			symbols[s.Symbol] = true
		}
	}
	return symbols, nil
}

// updateSymbolCache updates the cache with new and delisted symbols
func updateSymbolCache() (newSymbols, delistedSymbols []string) {
	mutex.Lock()
	defer mutex.Unlock()

	currentSymbols, err := fetchSymbolsFromBinance()
	if err != nil {
		log.Printf("Error fetching symbol data: %v", err)
		return nil, nil
	}

	for symbol := range currentSymbols {
		if !lastSymbols[symbol] {
			newSymbols = append(newSymbols, symbol)
		}
	}

	for symbol := range lastSymbols {
		if !currentSymbols[symbol] {
			delistedSymbols = append(delistedSymbols, symbol)
		}
	}

	lastSymbols = currentSymbols

	return newSymbols, delistedSymbols
}
