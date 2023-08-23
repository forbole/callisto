package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// handleMsgSetExtra allows to properly handle a handleMsgSetExtra
func (m *Module) handleMsgSetExtra(tx *juno.Tx, index int, msg *types.MsgSetExtra) error {
	account, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.Hash))
	switch {
	case err != nil:
		return err
	case len(account) != 1:
		return types.ErrInvalidHashField
	}

	account[0].Extras = msg.Extras

	if err := m.accountRepo.UpdateAccounts(account...); err != nil {
		return err
	}

	return m.accountRepo.SaveExtra(msg, tx.TxHash)
}
