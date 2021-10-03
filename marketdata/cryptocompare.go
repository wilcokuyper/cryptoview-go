package marketdata

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type CryptoClient interface {
	GetSymbols() (map[string]Symbol, error)
	GetPrice(symbol string, baseCurrency string) (float64, error)
}

type CryptocompareClient struct {
	key string
	baseURL string
	client *http.Client
}

func NewCryptocompareClient(key string, baseURL string) *CryptocompareClient {
	return &CryptocompareClient{key, baseURL, &http.Client{}}
}

type Symbol struct {
	Name string `json:"Name"`
	Symbol string `json:"Symbol"`
	CoinName string `json:"CoinName"`
}

func (c *CryptocompareClient) GetSymbols() (map[string]Symbol, error) {
	response, err := c.doRequest("GET", c.baseURL + "data/all/coinlist")
	if err != nil {
		return nil, errors.Wrap(err, "Unable to retrieve symbols")
	}

	var symbols struct{
		Data map[string]Symbol `json:"Data"`
	}

	err = json.Unmarshal(response, &symbols)

	if err != nil {
		return nil, err
	}

	return symbols.Data, nil
}

func (c *CryptocompareClient) GetPrice(symbol string, baseCurrency string) (float64, error) {
	response, err := c.doRequest("GET", c.baseURL + "data/price?fsym=" + symbol + "&tsyms=EUR")

	if err != nil {
		return -1, errors.Wrap(err, "Unable to perform request")
	}

	var parsed map[string]float64

	json.Unmarshal(response, &parsed)

	price, ok := parsed[baseCurrency]

	if !ok {
		return -1, errors.New("Price not found")
	}

	return price, nil
}

func (c *CryptocompareClient) doRequest(method string, url string) ([]byte, error) {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("authorization", "Apikey " + c.key)

	res, err := c.client.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "Cryptocompare error")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, errors.Wrap(err, "Cryptoview error")
	}

	return body, nil
}