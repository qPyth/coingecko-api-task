package coingecko

import "errors"

var (
	ErrCoinNotFound = errors.New("coin not found")
)

type CoinCurrency struct {
	Id           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	CurrentPrice float64 `json:"current_price"`
}
