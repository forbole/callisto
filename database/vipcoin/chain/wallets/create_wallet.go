package wallets

import (
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveCreateWallet saves the given wallets inside the database
func (r *Repository) SaveCreateWallet(msgWallet *walletstypes.MsgCreateWallet, transactionHash string) error {
	query := `INSERT INTO vipcoin_chain_wallets_create_wallet 
			(transaction_hash, creator, address, account_address, kind, state, extras) 
		VALUES 
			(:transaction_hash, :creator, :address, :account_address, :kind, :state, :extras)`

	if _, err := r.db.NamedExec(query, toCreateWalletDatabase(msgWallet, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetCreateWallet get the given wallet from database
func (r *Repository) GetCreateWallet(walletFilter filter.Filter) ([]*walletstypes.MsgCreateWallet, error) {
	query, args := walletFilter.Build(
		"vipcoin_chain_wallets_create_wallet",
		`creator, address, account_address, kind, state, extras`,
	)

	var result []types.DBCreateWallet

	if err := r.db.Select(&result, query, args...); err != nil {
		return []*walletstypes.MsgCreateWallet{}, errs.Internal{Cause: err.Error()}
	}

	wallets := make([]*walletstypes.MsgCreateWallet, 0, len(result))
	for _, wallet := range result {
		wallets = append(wallets, toCreateWalletDomain(wallet))
	}

	return wallets, nil
}
