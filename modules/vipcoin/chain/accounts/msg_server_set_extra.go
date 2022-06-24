package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSetExtra allows to properly handle a handleMsgSetExtra
func (m *Module) handleMsgSetExtra(tx *juno.Tx, index int, msg *types.MsgSetExtra) error {
	if err := m.accountRepo.SaveExtra(msg, tx.TxHash); err != nil {
		return err
	}

	account, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.Hash))
	if err != nil {
		return err
	}

	if len(account) != 1 {
		return types.ErrInvalidHashField
	}

	account[0].Extras = msg.Extras

	return m.accountRepo.UpdateAccounts(account...)
}
