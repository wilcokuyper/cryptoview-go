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
	res, err := w.db.Query("SELECT * FROM wallet_items WHERE user_id = ?", userId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetWallet: unable to fetch wallet for user: %v", userId)
	}
	defer res.Close()

	var items []*WalletItem

	err = res.Scan(items)
	if err != nil {
		return nil, errors.Wrap(err, "WalletItemsForUserId: unable scan rows into WalletItem struct")
	}

	return items, nil
}
