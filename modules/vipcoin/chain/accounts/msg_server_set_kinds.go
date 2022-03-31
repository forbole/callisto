package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"
)

// handleMsgSetKinds allows to properly handle a handleMsgSetKinds
func (m *Module) handleMsgSetKinds(tx *juno.Tx, index int, msg *types.MsgSetKinds) error {
	if err := m.accountRepo.SaveKinds(msg); err != nil {
		return err
	}

	acc, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(FieldHash, msg.Hash))
	if err != nil {
		return err
	}

	if len(acc) != 1 {
		return types.ErrInvalidHashField
	}

	acc[0].Kinds = msg.Kinds

	return m.accountRepo.UpdateAccounts(acc...)
}
