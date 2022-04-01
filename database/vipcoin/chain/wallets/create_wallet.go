package wallets

import (
	"context"
	"database/sql"

	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveCreateWallet saves the given wallets inside the database
func (r *Repository) SaveCreateWallet(msgWallet ...*walletstypes.MsgCreateWallet) error {
	if len(msgWallet) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO vipcoin_chain_wallets_create_wallet 
			(creator, address, account_address, kind, state, extras) 
		VALUES 
			(:creator, :address, :account_address, :kind, :state, :extras)`

	for _, wallet := range msgWallet {
		if _, err := tx.NamedExec(query, toCreateWalletDatabase(wallet)); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetCreateWallet get the given wallet from database
func (r *Repository) GetCreateWallet(walletFilter filter.Filter) ([]*walletstypes.MsgCreateWallet, error) {
	query, args := walletFilter.Build(
		"vipcoin_chain_wallets_create_wallet",
		`creator, address, account_address, kind, state, extras`,
	)

	var result []types.DBCreateWallet

	if err := r.db.Select(&result, query, args...); err != nil {
		return []*walletstypes.MsgCreateWallet{}, err
	}

	wallets := make([]*walletstypes.MsgCreateWallet, 0, len(result))
	for _, wallet := range result {
		wallets = append(wallets, toCreateWalletDomain(wallet))
	}

	return wallets, nil
}
