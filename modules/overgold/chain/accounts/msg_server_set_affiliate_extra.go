package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// handleMsgSetAffiliateExtra allows to properly handle a handleMsgSetAffiliateExtra
func (m *Module) handleMsgSetAffiliateExtra(tx *juno.Tx, index int, msg *types.MsgSetAffiliateExtra) error {
	acc, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.AccountHash))
	switch {
	case err != nil:
		return err
	case len(acc) != 1:
		return types.ErrInvalidHashField
	}

	affAcc, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.AffiliationHash))
	switch {
	case err != nil:
		return err
	case len(affAcc) != 1:
		return types.ErrInvalidHashField
	}

	for index, a := range acc[0].Affiliates {
		if a.Address == affAcc[0].Address {
			acc[0].Affiliates[index].Extras = msg.Extras
		}
	}

	if err := m.accountRepo.UpdateAccounts(acc...); err != nil {
		return err
	}

	return m.accountRepo.SaveAffiliateExtra(msg, tx.TxHash)
}
