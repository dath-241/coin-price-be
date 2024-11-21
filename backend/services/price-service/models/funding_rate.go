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
	Symbol                   string `json:"symbol" example:"QTUMUSDT"`
	FundingRate              string `json:"fundingRate" example:"0.00010000"`
	FundingCountDown         string `json:"fundingCountDown" example:"06:47:47"`
	EventTime                string `json:"eventTime" example:"2024-11-21 08:12:13"`
	AdjustedFundingRateCap   string `json:"adjustedFundingRateCap" example:"0.02000000"`
	AdjustedFundingRateFloor string `json:"adjustedFundingRateFloor" example:"-0.02000000"`
	FundingIntervalHours     int    `json:"fundingIntervalHours" example:"8"`
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

type FundingRateWebSocket struct {
	Data struct {
		Symbol          string `json:"s"`
		E               string `json:"e"`
		EventTime       int64  `json:"E"`
		FundingRate     string `json:"r"`
		NextFundingTime int64  `json:"T"`
	} `json:"data"`
}
