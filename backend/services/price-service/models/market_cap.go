package models

type CategoriesInput struct {
	Id    string `json:"id"`
	Start int64  `json:"start"`
	Limit int64  `json:"limit"`
}

type CategoryInput struct {
	IdCategory string `json:"id_category"`
	Start      int64  `json:"start"`
	Limit      int64  `json:"limit"`
	Convert    string `json:"convert"`
}
type CoinInput struct {
	IdCoin    int64  `json:"id_coin"`
	ConvertId string `json:"convert_id"`
}
