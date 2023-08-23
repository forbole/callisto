package wallets

import (
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v3/database/types"
)

// SaveCreateWalletWithBalance - saves the given wallet inside the database
func (r Repository) SaveCreateWalletWithBalance(msg *walletstypes.MsgCreateWalletWithBalance, transactionHash string) error {
	query := `INSERT INTO overgold_chain_wallets_create_wallet_with_balance 
			(transaction_hash, creator, address, account_address, kind, state, extras, default_status, balance) 
		VALUES 
			(:transaction_hash, :creator, :address, :account_address, :kind, :state, :extras, :default_status, :balance)`

	if _, err := r.db.NamedExec(query, toCreateWalletWithBalanceDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetCreateWalletWithBalance - get the given wallet from database
func (r Repository) GetCreateWalletWithBalance(walletFilter filter.Filter) ([]*walletstypes.MsgCreateWalletWithBalance, error) {
	query, args := walletFilter.Build(
		"overgold_chain_wallets_create_wallet_with_balance",
		types.FieldCreator, types.FieldAddress, types.FieldAccountAddress, types.FieldKind, types.FieldState,
		types.FieldExtras, types.FieldDefaultStatus, types.FieldBalance,
	)

	var result []types.DBCreateWalletWithBalance
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*walletstypes.MsgCreateWalletWithBalance{}, errs.Internal{Cause: err.Error()}
	}

	wallets := make([]*walletstypes.MsgCreateWalletWithBalance, 0, len(result))
	for _, wallet := range result {
		wallets = append(wallets, toCreateWalletWithBalanceDomain(wallet))
	}

	return wallets, nil
}
