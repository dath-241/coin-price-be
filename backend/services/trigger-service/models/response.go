package models

// General Error Response Types
type ErrorResponse struct {
	Message string `json:"message" example:"An error occurred"`
}

// Alert Response Types
type ResponseAlertCreated struct {
	Message string `json:"message" example:"Alert created successfully"`
	AlertID string `json:"alert_id" example:"647f1f77bcf86cd799439011"`
}

type ResponseAlertDetail struct {
	ID              string    `json:"id" example:"647f1f77bcf86cd799439011"`
	Symbol          string    `json:"symbol" example:"BTCUSDT"`
	Condition       string    `json:"condition" example:">="`
	Price           float64   `json:"price" example:"50000"`
	IsActive        bool      `json:"isActive" example:"true"`
	CreatedAt       string    `json:"createdAt" example:"2024-11-21T00:00:00Z"`
	UpdatedAt       string    `json:"updatedAt" example:"2024-11-21T00:00:00Z"`
	Frequency       string    `json:"frequency" example:"immediate"`
	MaxRepeatCount  int       `json:"maxRepeatCount" example:"5"`
	SnoozeCondition string    `json:"snoozeCondition" example:"none"`
	Range           []float64 `json:"range"`
}

type ResponseAlertList struct {
	Alerts []ResponseAlertDetail `json:"alerts"`
}

type ResponseAlertDeleted struct {
	Message string `json:"message" example:"Alert deleted successfully"`
}

// Symbol Response Types
type ResponseSymbolUpdate struct {
	NewSymbols      []string `json:"new_symbols" example:"[BTCUSDT, ETHUSDT]"`
	DelistedSymbols []string `json:"delisted_symbols" example:"[BNBUSDT]"`
}

// Price and Funding Data Response Types
type ResponsePrice struct {
	Symbol string  `json:"symbol" example:"BTCUSDT"`
	Price  float64 `json:"price" example:"50000.75"`
}

type ResponseFundingRate struct {
	Symbol       string  `json:"symbol" example:"BTCUSDT"`
	FundingRate  float64 `json:"funding_rate" example:"0.0001"`
	IntervalTime string  `json:"interval_time" example:"8h"`
}

type ResponsePriceDifference struct {
	Symbol          string  `json:"symbol" example:"BTCUSDT"`
	SpotPrice       float64 `json:"spot_price" example:"50000.75"`
	FuturePrice     float64 `json:"future_price" example:"50200.50"`
	PriceDifference float64 `json:"price_difference" example:"199.75"`
}

// New and Delisted Symbols Response
type ResponseNewDelistedSymbols struct {
	NewSymbols      []string `json:"new_symbols" example:"[BTCUSDT, ETHUSDT]"`
	DelistedSymbols []string `json:"delisted_symbols" example:"[BNBUSDT]"`
}

// Set Symbol Alert Response
type ResponseSetSymbolAlert struct {
	Message string `json:"message" example:"Alert created successfully"`
	AlertID string `json:"alert_id" example:"647f1f77bcf86cd799439011"`
}

type ResponseAlertCheckerStatus struct {
	Status string `json:"status" example:"Alert checker started"`
}

// Response for CreateUser success
type ResponseUserCreated struct {
	Message string `json:"message" example:"User created successfully"`
	UserID  string `json:"user_id" example:"647f1f77bcf86cd799439011"`
}

// Response for GetUserAlerts success
type ResponseUserAlerts struct {
	Alerts []Alert `json:"alerts"`
}

// Response for NotifyUser success
type ResponseNotificationSent struct {
	Status string `json:"status" example:"Notification sent"`
}
