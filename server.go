package main

import (
	"net/http"
	"os"

	"github.com/wilcokuyper/cryptoview-go/marketdata"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	mux *http.ServeMux
	client marketdata.CryptoClient
}

func NewServer(logger *zap.Logger, mux *http.ServeMux, client marketdata.CryptoClient) *Server {
	return &Server{
		logger: logger,
		mux: mux,
		client: client,
	}
}

func (s *Server) run() {
	// Setup API routes
	s.setupRoutes()

	// Lookup port and start server
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	logger.Info("Starting webserver. Listening on :" + port)

	err := http.ListenAndServe(":" + port, s.mux)
	if err != nil {
		logger.Fatal("Unable to start server", zap.Error(err))
	}
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/api/price", s.GetPriceHandler)
	s.mux.HandleFunc("/api/symbols", s.GetSymbolsHandler)
	s.mux.HandleFunc("/api/historical-data", s.GetHistoricalDataHandler)
}