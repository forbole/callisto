package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSetAffiliateExtra allows to properly handle a handleMsgSetAffiliateExtra
func (m *Module) handleMsgSetAffiliateExtra(tx *juno.Tx, index int, msg *types.MsgSetAffiliateExtra) error {
	if err := m.accountRepo.SaveAffiliateExtra(msg); err != nil {
		return err
	}

	acc, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.AccountHash))
	if err != nil {
		return err
	}

	affAcc, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.AffiliationHash))
	if err != nil {
		return err
	}

	if len(acc) != 1 || len(affAcc) != 1 {
		return types.ErrInvalidHashField
	}

	for index, a := range acc[0].Affiliates {
		if a.Address == affAcc[0].Address {
			acc[0].Affiliates[index].Extras = msg.Extras
		}
	}

	return m.accountRepo.UpdateAccounts(acc...)
}
