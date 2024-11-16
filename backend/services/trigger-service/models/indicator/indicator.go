package indicator

type Indicator struct {
	ID                 string `bson:"_id,omitempty" json:"id"`
	Symbol             string `json:"symbol" binding:"required"`
	Indicator          string `json:"indicator" binding:"required"`
	Period             int    `json:"period" binding:"required"`
	NotificationMethod string `json:"notification_method" binding:"required"`
}
