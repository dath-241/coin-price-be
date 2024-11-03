package models

type FundingRateFirst struct {
	Symbol          string `json:"symbol"`
	FundingRate     string `json:"lastFundingRate"`
	NextFundingTime int64  `json:"nextFundingTime"`
	EventTime       int64  `json:"time"`
}

type FundingRateSecond struct {
	Symbol                   string `json:"symbol"`
	AdjustedFundingRateCap   string `json:"adjustedFundingRateCap"`
	AdjustedFundingRateFloor string `json:"adjustedFundingRateFloor"`
	FundingIntervalHours     int    `json:"fundingIntervalHours"`
}

type ResponseFundingRate struct {
	Symbol                   string `json:"symbol"`
	FundingRate              string `json:"fundingRate"`
	FundingCountDown         string `json:"fundingCountDown"`
	EventTime                string `json:"eventTime"`
	AdjustedFundingRateCap   string `json:"adjustedFundingRateCap"`
	AdjustedFundingRateFloor string `json:"adjustedFundingRateFloor"`
	FundingIntervalHours     int    `json:"fundingIntervalHours"`
}

func (r *ResponseFundingRate) UpdateData(symbol, fundingRate, fundingCountDown, eventTime, adjustedFundingRateCap, AdjustedFundingRateFloor string, fundingIntervalHours int) {
	r.Symbol = symbol
	r.FundingRate = fundingRate
	r.FundingCountDown = fundingCountDown
	r.EventTime = eventTime
	r.AdjustedFundingRateCap = adjustedFundingRateCap
	r.AdjustedFundingRateFloor = AdjustedFundingRateFloor
	r.FundingIntervalHours = fundingIntervalHours
}
