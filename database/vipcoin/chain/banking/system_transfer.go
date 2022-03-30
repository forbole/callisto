/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"
)

// SaveSystemTransfers - method that create system transfers to the "vipcoin_chain_banking_system_transfer" table
func (r Repository) SaveSystemTransfers(transfers ...*bankingtypes.SystemTransfer) error {
	return nil
}

// GetSystemTransfers - method that get system transfers from the "vipcoin_chain_banking_system_transfer" table
func (r Repository) GetSystemTransfers(filter filter.Filter) ([]*bankingtypes.SystemTransfer, error) {
	return []*bankingtypes.SystemTransfer{}, nil
}

// UpdateSystemTransfers - method that update the transfer in the "vipcoin_chain_banking_system_transfer" table
func (r Repository) UpdateSystemTransfers(transfers ...*bankingtypes.SystemTransfer) error {
	return nil
}
