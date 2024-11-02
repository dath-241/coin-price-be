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
