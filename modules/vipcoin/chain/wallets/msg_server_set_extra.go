package wallets

import (
	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSetExtra allows to properly handle a MsgSetExtra
func (m *Module) handleMsgSetExtra(msg *typeswallets.MsgSetExtra) error {
	if err := m.walletsRepo.SaveExtras(msg); err != nil {
		return err
	}

	wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.Address))
	switch {
	case err != nil:
		return err
	case len(wallets) != 1:
		return typeswallets.ErrInvalidAddressField
	}

	wallets[0].Extras = msg.Extras

	return m.walletsRepo.SaveWallets(wallets...)
}
