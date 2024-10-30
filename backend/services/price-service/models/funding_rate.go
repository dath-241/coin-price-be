package models

// use for get realtime funding rate
type InputFundingRate struct {
	Symbol string `json:"symbol"`
}

// use for get funding rate countdown
type ResponseFundingRate struct {
	Symbol          string `json:"symbol"`
	NextFundingTime int64  `json:"nextFundingTime"`
	CurrentTime     int64  `json:"time"`
}
