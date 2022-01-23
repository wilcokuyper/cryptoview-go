package marketdata

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type MarketdataHandler struct {
	logger *zap.Logger
	client CryptoClient
}

func NewMarketdataHandler(logger *zap.Logger, client CryptoClient) *MarketdataHandler {
	return &MarketdataHandler{
		logger: logger,
		client: client,
	}
}

func (s *MarketdataHandler) HandlePrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("GetPrice", zap.Any("query", r.URL.Query()))
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if !r.URL.Query().Has("symbol") {
			w.WriteHeader(http.StatusNotAcceptable)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Missing symbol parameter",
			})
			if err != nil {
				return
			}
			return
		}

		symbol := r.URL.Query().Get("symbol")
		price, err := s.client.GetPrice(symbol, "EUR")

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Price not found",
			})
			if err != nil {
				return
			}

			return
		}

		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"price": map[string]interface{}{
				"symbol": symbol,
				"price":  price,
			},
		})
		if err != nil {
			return
		}
	}
}

func (s *MarketdataHandler) HandleSymbols() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("GetSymbols")
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		symbols, err := s.client.GetSymbols()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Unable to retrieve symbols",
			})
			if err != nil {
				return
			}

			return
		}

		err = json.NewEncoder(w).Encode(symbols)
		if err != nil {
			return
		}
	}
}

func (s *MarketdataHandler) HandleHistoricalData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("GetHistoricalData", zap.Any("query", r.URL.Query()))
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if !r.URL.Query().Has("symbol") {
			w.WriteHeader(http.StatusNotAcceptable)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Missing symbol parameter",
			})
			if err != nil {
				return
			}
			return
		}

		symbol := r.URL.Query().Get("symbol")
		dataPoints, err := s.client.GetHistoricalData(symbol, "EUR", 10)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			err = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Price not found",
			})
			if err != nil {
				return
			}

			return
		}

		err = json.NewEncoder(w).Encode(dataPoints)
		if err != nil {
			return
		}
	}
}

func (s *MarketdataHandler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/price", s.HandlePrice())
	mux.HandleFunc("/api/symbols", s.HandleSymbols())
	mux.HandleFunc("/api/historical-data", s.HandleHistoricalData())
}
