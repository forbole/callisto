package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"
)

// SaveWithdraws - method that create withdraws to the "vipcoin_chain_banking_withdraw" table
func (r Repository) SaveWithdraws(withdraws ...*bankingtypes.Withdraw) error {
	return nil
}

// GetWithdraws - method that get withdraws from the "vipcoin_chain_banking_withdraw" table
func (r Repository) GetWithdraws(filter filter.Filter) ([]*bankingtypes.Withdraw, error) {
	return []*bankingtypes.Withdraw{}, nil
}

// UpdateWithdraws - method that update the withdraw in the "vipcoin_chain_banking_withdraw" table
func (r Repository) UpdateWithdraws(withdraws ...*bankingtypes.Withdraw) error {
	return nil
}
