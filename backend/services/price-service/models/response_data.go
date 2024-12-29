package models

type ErrorResponseDataMissing struct {
	Message string `json:"message" example:"Invalid symbol or request parameters"`
}
type ErrorResponseDataNotFound struct {
	Message string `json:"message" example:"Symbol not found"`
}
type ErrorResponseDataInternalServerError struct {
	Message string `json:"message" example:"Internal server error"`
}

type ErrorResponseInputMissing struct {
	Message string `json:"message" example:"Missing data"`
}

type ResponseKline struct {
	Symbol    string           `json:"symbol" example:"BTCUSDT"`
	Interval  string           `json:"interval" example:"1m"`
	EventTime string           `json:"eventTime" example:"2024-11-21 08:37:58"`
	KlineData []KlineDataPoint `json:"kline_data"`
}

type KlineDataPoint struct {
	Time   string  `json:"time" example:"2024-11-21T00:00:00Z"`
	Open   float64 `json:"open" example:"94288"`
	High   float64 `json:"high" example:"95000"`
	Low    float64 `json:"low" example:"94200"`
	Close  float64 `json:"close" example:"94645.7"`
	Volume float64 `json:"volume" example:"21752.097"`
}
