package models

type FutureKlineWebSocket struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	Kline     struct {
		StartTime    int64  `json:"t"`
		CloseTime    int64  `json:"T"`
		Symbol       string `json:"s"`
		Interval     string `json:"i"`
		FirstTradeID int64  `json:"f"`
		LastTradeID  int64  `json:"L"`
		OpenPrice    string `json:"o"`
		ClosePrice   string `json:"c"`
		HighPrice    string `json:"h"`
		LowPrice     string `json:"l"`
		Volume       string `json:"v"`
		TradeCount   int    `json:"n"`
		IsFinal      bool   `json:"x"`
		QuoteVolume  string `json:"q"`
		ActiveVolume string `json:"V"`
		ActiveQty    string `json:"Q"`
		BidVolume    string `json:"B"`
	} `json:"k"`
}

type ResponseFuturePrice struct {
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	EventTime string `json:"eventTime"`
}

func (r *ResponseFuturePrice) UpdateData(symbol, price, eventTime string) {
	r.Symbol = symbol
	r.Price = price
	r.EventTime = eventTime
}
