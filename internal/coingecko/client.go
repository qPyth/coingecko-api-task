package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type GCClient struct {
	host     string
	basePath string
	client   http.Client
	token    string
}

func NewCGClient(host, path string, token string) *GCClient {
	return &GCClient{host: host, basePath: path, client: http.Client{}, token: token}
}

func (c *GCClient) FetchAllCurrencies() ([]CoinCurrency, error) {

	params := url.Values{}

	params.Set("vs_currency", "usd")
	params.Set("order", "market_cap_desc")
	params.Set("per_page", "250")
	params.Set("page", "1")
	if c.token != "" {
		params.Set("x_cg_demo_api_key", c.token)
	}

	data, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var coins []CoinCurrency

	err = json.Unmarshal(data, &coins)
	if err != nil {
		return nil, err
	}

	return coins, err
}

func (c *GCClient) doRequest(params url.Values) ([]byte, error) {
	op := "client.doRequest"

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   c.basePath,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s new request error: %w", op, err)
	}
	req.URL.RawQuery = params.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s do request error: %w", op, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s reading request body error: %w", op, err)
	}

	return body, nil
}
