package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSetAffiliateAddress allows to properly handle a handleMsgSetAffiliateAddress
func (m *Module) handleMsgSetAffiliateAddress(tx *juno.Tx, index int, msg *types.MsgSetAffiliateAddress) error {
	acc, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.Hash))
	switch {
	case err != nil:
		return err
	case len(acc) != 1:
		return types.ErrInvalidHashField
	}

	if err := updateAffiliateAddress(acc[0].Affiliates, msg); err != nil {
		return err
	}

	if err := m.accountRepo.UpdateAccounts(acc...); err != nil {
		return err
	}

	return m.accountRepo.SaveAffiliateAddress(msg, tx.TxHash)
}
