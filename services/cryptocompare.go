package services

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Cryptocompare struct {
	key string
	baseURL string
	client *http.Client
}

func NewCryptocompare(key string, baseURL string) *Cryptocompare {
	return &Cryptocompare{key, baseURL, &http.Client{}}
}

type Symbol struct {
	Name string `json:"Name"`
	Symbol string `json:"Symbol"`
	CoinName string `json:"CoinName"`
}

func (c *Cryptocompare) GetSymbols() (map[string]Symbol, error) {
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

func (c *Cryptocompare) GetPrice(symbol string, baseCurrency string) (float64, error) {
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

func (c *Cryptocompare) doRequest(method string, url string) ([]byte, error) {
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