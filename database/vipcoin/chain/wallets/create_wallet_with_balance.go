package wallets

import (
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"

	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveCreateWalletWithBalance - saves the given wallet inside the database
func (r Repository) SaveCreateWalletWithBalance(msg ...*walletstypes.MsgCreateWalletWithBalance) error {
	if len(msg) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_wallets_create_wallet_with_balance 
			(creator, address, account_address, kind, state, extras, default_status, balance) 
		VALUES 
			(:creator, :address, :account_address, :kind, :state, :extras, :default_status, :balance)`

	if _, err := r.db.NamedExec(query, toCreateWalletsWithBalanceDatabase(msg...)); err != nil {
		return err
	}

	return nil
}

// GetCreateWalletWithBalance - get the given wallet from database
func (r Repository) GetCreateWalletWithBalance(walletFilter filter.Filter) ([]*walletstypes.MsgCreateWalletWithBalance, error) {
	query, args := walletFilter.Build(
		"vipcoin_chain_wallets_create_wallet_with_balance",
		types.FieldCreator, types.FieldAddress, types.FieldAccountAddress, types.FieldKind, types.FieldState,
		types.FieldExtras, types.FieldDefaultStatus, types.FieldBalance,
	)

	var result []types.DBCreateWalletWithBalance
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*walletstypes.MsgCreateWalletWithBalance{}, err
	}

	wallets := make([]*walletstypes.MsgCreateWalletWithBalance, 0, len(result))
	for _, wallet := range result {
		wallets = append(wallets, toCreateWalletWithBalanceDomain(wallet))
	}

	return wallets, nil
}
