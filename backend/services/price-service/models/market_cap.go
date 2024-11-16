package models

type MarketCapResponse struct {
	Symbol     string `json:"symbol"`
	MarketData struct {
		MarketCap struct {
			USD int64 `json:"usd"`
		} `json:"market_cap"`
		TotalVolume struct {
			USD int64 `json:"usd"`
		} `json:"total_volume"`
	} `json:"market_data"`
}

type FormatMarketCapResponse struct {
	Symbol      string `json:"symbol"`
	MarketCap   int64  `json:"market_cap"`
	TotalVolume int64  `json:"24h_volume"`
}

func CreateReponseFormat(symbol string, marketCap, totalVolume int64) *FormatMarketCapResponse {
	return &FormatMarketCapResponse{
		Symbol:      symbol,
		MarketCap:   marketCap,
		TotalVolume: totalVolume,
	}
}
