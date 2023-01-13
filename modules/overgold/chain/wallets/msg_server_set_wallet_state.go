package wallets

import (
	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSetStates allows to properly handle a MsgSetState
func (m *Module) handleMsgSetStates(tx *juno.Tx, index int, msg *typeswallets.MsgSetWalletState) error {
	wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.Address))
	switch {
	case err != nil:
		return err
	case len(wallets) != 1:
		return typeswallets.ErrInvalidAddressField
	}

	wallets[0].State = msg.State

	if err := m.walletsRepo.UpdateWallets(wallets...); err != nil {
		return err
	}

	return m.walletsRepo.SaveStates(msg, tx.TxHash)
}
