package main

import (
	"database/sql"
	"github.com/wilcokuyper/cryptoview-go/services/marketdata"
	"github.com/wilcokuyper/cryptoview-go/services/wallet"
	"go.uber.org/zap"
	"net/http"
	"os"
)

type server struct {
	logger *zap.Logger
	db     *sql.DB

	mux *http.ServeMux
}

func NewServer(logger *zap.Logger, db *sql.DB) *server {
	mux := http.NewServeMux()
	return &server{logger, db, mux}
}

func (s *server) Run(port string) {
	logger.Info("Starting webserver. Listening on :" + port)

	s.setupMarketdataHandler()
	s.setupWalletHandler()

	err := http.ListenAndServe(":"+port, s.mux)
	if err != nil {
		logger.Fatal("Unable to start server", zap.Error(err))
	}
}

func (s *server) setupMarketdataHandler() {
	client := marketdata.NewCryptocompareClient(
		os.Getenv("CRYPTOCOMPARE_API_KEY"),
		os.Getenv("CRYPTOCOMPARE_BASE_URL"),
		logger,
	)

	marketdataHandler := marketdata.NewMarketdataHandler(logger, client)
	marketdataHandler.SetupRoutes(s.mux)
}

func (s *server) setupWalletHandler() {
	client := wallet.NewWalletRepository(s.db)

	walletHandler := wallet.NewWalletHandler(s.logger, client)
	walletHandler.SetupRoutes(s.mux)
}
