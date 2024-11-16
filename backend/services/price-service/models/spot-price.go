package models

type SpotTickerWebSocket struct {
	EventType        string `json:"e"`
	EventTime        int64  `json:"E"`
	Symbol           string `json:"s"`
	PriceChange      string `json:"p"`
	PriceChangePct   string `json:"P"`
	WeightedAvgPrice string `json:"w"`
	PrevClosePrice   string `json:"x"`
	LastPrice        string `json:"c"`
	LastQty          string `json:"Q"`
	BidPrice         string `json:"b"`
	BidQty           string `json:"B"`
	AskPrice         string `json:"a"`
	AskQty           string `json:"A"`
	OpenPrice        string `json:"o"`
	HighPrice        string `json:"h"`
	LowPrice         string `json:"l"`
	Volume           string `json:"v"`
	QuoteVolume      string `json:"q"`
	OpenTime         int64  `json:"O"`
	CloseTime        int64  `json:"C"`
	FirstTradeID     int64  `json:"F"`
	LastTradeID      int64  `json:"L"`
	TradeCount       int    `json:"n"`
}

type ResponseSpotPrice struct {
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	EventTime string `json:"eventTime"`
}

func (r *ResponseSpotPrice) UpdateData(symbol, price, eventTime string) {
	r.Symbol = symbol
	r.Price = price
	r.EventTime = eventTime
}
