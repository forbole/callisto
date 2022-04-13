package wallets

import (
	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSetStates allows to properly handle a MsgSetState
func (m *Module) handleMsgSetStates(msg *typeswallets.MsgSetWalletState) error {
	if err := m.walletsRepo.SaveStates(msg); err != nil {
		return err
	}

	wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.Address))
	switch {
	case err != nil:
		return err
	case len(wallets) != 1:
		return typeswallets.ErrInvalidAddressField
	}

	wallets[0].State = msg.State

	return m.walletsRepo.SaveWallets(wallets...)
}
