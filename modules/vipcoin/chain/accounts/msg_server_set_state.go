package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSetState allows to properly handle a handleMsgSetState
func (m *Module) handleMsgSetState(tx *juno.Tx, index int, msg *types.MsgSetState) error {
	acc, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.Hash))
	switch {
	case err != nil:
		return err
	case len(acc) != 1:
		return types.ErrInvalidHashField
	}

	acc[0].State = msg.State

	if err := m.accountRepo.UpdateAccounts(acc...); err != nil {
		return err
	}

	return m.accountRepo.SaveState(msg, tx.TxHash)
}
