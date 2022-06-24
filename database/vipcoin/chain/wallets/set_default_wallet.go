package wallets

import (
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveDefaultWallets - inserts messages into the "vipcoin_chain_wallets_set_default_wallet" table
func (r Repository) SaveDefaultWallets(messages *walletstypes.MsgSetDefaultWallet, transactionHash string) error {
	query := `INSERT INTO vipcoin_chain_wallets_set_default_wallet 
			(transaction_hash, creator, address) 
			VALUES 
			(:transaction_hash, :creator, :address)`

	if _, err := r.db.NamedExec(query, toSetDefaultWalletDatabase(messages, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetDefaultWallets - get the given messages from the "vipcoin_chain_wallets_set_default_wallet" table
func (r Repository) GetDefaultWallets(filter filter.Filter) ([]*walletstypes.MsgSetDefaultWallet, error) {
	query, args := filter.Build("vipcoin_chain_wallets_set_default_wallet", types.FieldCreator, types.FieldAddress)

	var wallets []types.DBSetDefaultWallet
	if err := r.db.Select(&wallets, query, args...); err != nil {
		return []*walletstypes.MsgSetDefaultWallet{}, errs.Internal{Cause: err.Error()}
	}

	result := make([]*walletstypes.MsgSetDefaultWallet, 0, len(wallets))
	for _, w := range wallets {
		result = append(result, toSetDefaultWalletDomain(w))
	}

	return result, nil
}
