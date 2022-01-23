package marketdata

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Symbol struct {
	Name     string `json:"Name"`
	Symbol   string `json:"Symbol"`
	CoinName string `json:"CoinName"`
}

type DataPoint struct {
	Time       int64   `json:"time"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Open       float64 `json:"open"`
	VolumeFrom float64 `json:"volumefrom"`
	VolumeTo   float64 `json:"volumeto"`
	Close      float64 `json:"close"`
}

type CryptoClient interface {
	GetSymbols() (map[string]Symbol, error)
	GetPrice(symbol string, baseCurrency string) (float64, error)
	GetHistoricalData(symbols string, baseCurrency string, limit int) ([]DataPoint, error)
}

type CryptocompareClient struct {
	key     string
	baseURL string
	client  *http.Client
	logger  *zap.Logger
}

func NewCryptocompareClient(key string, baseURL string, logger *zap.Logger) *CryptocompareClient {
	return &CryptocompareClient{
		key,
		baseURL,
		&http.Client{Timeout: 10 * time.Second},
		logger,
	}
}

func (c *CryptocompareClient) GetSymbols() (map[string]Symbol, error) {
	c.logger.Info("GetSymbols")
	response, err := c.doRequest("GET", c.baseURL+"data/all/coinlist")
	if err != nil {
		return nil, errors.Wrap(err, "Unable to retrieve symbols")
	}

	var symbols struct {
		Data map[string]Symbol `json:"Data"`
	}

	err = json.Unmarshal(response, &symbols)
	if err != nil {
		return nil, err
	}

	return symbols.Data, nil
}

func (c *CryptocompareClient) GetPrice(symbol string, baseCurrency string) (float64, error) {
	c.logger.Info(
		"GetPrice",
		zap.String("symbol", symbol),
		zap.String("currency", baseCurrency),
	)
	response, err := c.doRequest("GET", c.baseURL+"data/price?fsym="+symbol+"&tsyms="+baseCurrency)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to perform get price for symbol "+symbol)
	}

	var parsed map[string]float64

	err = json.Unmarshal(response, &parsed)
	if err != nil {
		return 0, err
	}

	price, ok := parsed[baseCurrency]
	if !ok {
		return -1, errors.New("Price not found")
	}

	return price, nil
}

func (c *CryptocompareClient) GetHistoricalData(symbol string, baseCurrency string, limit int) ([]DataPoint, error) {
	c.logger.Info(
		"GetHistoricalData",
		zap.String("symbol", symbol),
		zap.String("currency", baseCurrency),
		zap.Int("limit", 16),
	)
	response, err := c.doRequest("GET", c.baseURL+"data/v2/histominute?fsym="+symbol+"&tsym="+baseCurrency+"&limit="+strconv.Itoa(limit))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get historical data for symbol "+symbol)
	}

	var data struct {
		Data struct {
			Data []DataPoint `json:"Data"`
		} `json:"Data"`
	}

	err = json.Unmarshal(response, &data)
	if err != nil {
		return nil, errors.New("No historical data found for symbol " + symbol)
	}

	return data.Data.Data, nil
}

func (c *CryptocompareClient) doRequest(method string, url string) ([]byte, error) {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("authorization", "Apikey "+c.key)

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
