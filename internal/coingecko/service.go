package coingecko

import (
	"log/slog"
	"strings"
	"sync"
)

type CurrencyProvider interface {
	FetchAllCurrencies() ([]CoinCurrency, error)
}

type CoinService struct {
	cache map[string]CoinCurrency
	mu    sync.Mutex
	cp    CurrencyProvider
	log   *slog.Logger
}

func NewCoinService(cp CurrencyProvider, log *slog.Logger) *CoinService {
	c := &CoinService{cache: make(map[string]CoinCurrency), cp: cp, log: log}
	c.AsyncUpdate()
	return c
}

func (c *CoinService) AsyncUpdate() {

	coins, err := c.cp.FetchAllCurrencies()
	switch err {
	case nil:
		c.mu.Lock()
		defer c.mu.Unlock()
		for _, coin := range coins {
			c.cache[strings.ToLower(coin.Name)] = coin
		}
		c.log.Info("coins updated successfully")
	default:
		c.log.Warn("coins update error", "err", err.Error())
	}
}

func (c *CoinService) CoinCurrency(name string) (CoinCurrency, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if coin, ok := c.cache[name]; ok {
		return coin, nil
	}

	c.log.Info("coin not found", "coin", name)
	return CoinCurrency{}, ErrCoinNotFound
}
