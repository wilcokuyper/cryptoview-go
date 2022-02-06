package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wilcokuyper/cryptoview-go/services/auth"
	"go.uber.org/zap"
)

type WalletItem struct {
	UserId   int
	Currency string
	Amount   float64
}

type walletClient interface {
	WalletItemsForUserId(userId int) ([]*WalletItem, error)
}

type WalletHandler struct {
	logger *zap.Logger
	client walletClient
}

func NewWalletHandler(logger *zap.Logger, client walletClient) *WalletHandler {
	return &WalletHandler{
		logger,
		client,
	}
}

func (s *WalletHandler) viewWallet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("viewWallet")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		user, _ := r.Context().Value(auth.ContextUser("user")).(auth.User)
		fmt.Printf("user: %v\n", user)
		items, err := s.client.WalletItemsForUserId(int(user.Id))
		if err != nil {
			s.logger.Info("viewWallet:", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(items)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s *WalletHandler) addItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("addItem")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *WalletHandler) updateItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("updateItem")

		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *WalletHandler) deleteItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("deleteWalletItem")

		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *WalletHandler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/wallet/view", auth.Middleware(s.viewWallet()))
	mux.HandleFunc("/wallet/add", auth.Middleware(s.addItem()))
	mux.HandleFunc("/wallet/update", auth.Middleware(s.updateItem()))
	mux.HandleFunc("/wallet/delete", auth.Middleware(s.deleteItem()))
}
