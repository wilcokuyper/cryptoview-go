package main

import (
	"encoding/json"
	"net/http"

	"github.com/wilcokuyper/cryptoview-go/marketdata"
)

type Server struct {
	client marketdata.CryptoClient
}

func (s *Server) GetPriceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !r.URL.Query().Has("symbol") {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Missing symbol parameter",
		})
		return
	}

	symbol := r.URL.Query().Get("symbol")
	price, err := s.client.GetPrice(symbol, "EUR")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Price not found",
		})

		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"price": map[string]interface{}{
			"symbol": symbol,
			"price": price,
		},
	})
}

func (s *Server) GetSymbolsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	symbols, err := s.client.GetSymbols()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Unable to retrieve symbols",
		})

		return
	}

	json.NewEncoder(w).Encode(symbols)
}

func (s *Server) GetHistoricalDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !r.URL.Query().Has("symbol") {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Missing symbol parameter",
		})
		return
	}

	symbol := r.URL.Query().Get("symbol")
	dataPoints, err := s.client.GetHistoricalData(symbol, "EUR", 10)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Price not found",
		})

		return
	}

	json.NewEncoder(w).Encode(dataPoints)
}