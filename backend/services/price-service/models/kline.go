package models

type KlineResponse struct {
	Symbol    string          `json:"symbol"`
	Interval  string          `json:"interval"`
	EventTime string          `json:"eventTime"`
	KlineData []KLineEachData `json:"kline_data"`
}

func (kline *KlineResponse) UpdateKlineResponse(symbol, interval, eventTime string) {
	kline.Symbol = symbol
	kline.Interval = interval
	kline.EventTime = eventTime
}
func (kline *KlineResponse) UpdateKlineResponseData(klineData *KLineEachData) {
	kline.KlineData = append(kline.KlineData, *klineData)
}

type KLineEachData struct {
	Time   string  `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

func (klineEach *KLineEachData) UpdateKlineEachData(time string, open, high, low, close, volume float64) {
	klineEach.Time = time
	klineEach.Open = open
	klineEach.High = high
	klineEach.Low = low
	klineEach.Close = close
	klineEach.Volume = volume
}

// struct for KlineWebsocket
type KlineWebsocket struct {
	Data struct {
		EventType string `json:"e"`
		EventTime int64  `json:"E"`
		Symbol    string `json:"s"`
		KData     struct {
			StartTime           int64  `json:"t"`
			CloseTime           int64  `json:"T"`
			LastTrade           int64  `json:"L"`
			OpenPrice           string `json:"o"`
			ClosePrice          string `json:"c"`
			HighPrice           string `json:"h"`
			LowPrice            string `json:"l"`
			BaseAssetVolume     string `json:"v"`
			QuoteAssetVolume    string `json:"q"`
			TakerBuyBaseVolume  string `json:"V"`
			TakerBuyQuoteVolume string `json:"Q"`
		} `json:"k"`
	} `json:"data"`
}
