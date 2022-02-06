package wallet

import (
	"database/sql"

	"github.com/pkg/errors"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db}
}

func (w *WalletRepository) WalletItemsForUserId(userId int) ([]*WalletItem, error) {
	q, err := w.db.Query("SELECT user_id, currency, amount FROM wallet_items WHERE user_id = ?", userId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetWallet: unable to fetch wallet for user: %v", userId)
	}
	defer q.Close()

	var items []*WalletItem

	for q.Next() {
		var item WalletItem
		if err = q.Scan(&item.UserId, &item.Currency, &item.Amount); err != nil {
			return nil, errors.Wrap(err, "WalletItemsForUserId: unable scan rows into WalletItem struct")
		}
		items = append(items, &item)
	}


	return items, nil
}
