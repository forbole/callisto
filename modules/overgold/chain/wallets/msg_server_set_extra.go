package wallets

import (
	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// handleMsgSetExtra allows to properly handle a MsgSetExtra
func (m *Module) handleMsgSetExtra(tx *juno.Tx, index int, msg *typeswallets.MsgSetExtra) error {

	wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.Address))
	switch {
	case err != nil:
		return err
	case len(wallets) != 1:
		return typeswallets.ErrInvalidAddressField
	}

	wallets[0].Extras = msg.Extras

	if err := m.walletsRepo.UpdateWallets(wallets...); err != nil {
		return err
	}

	return m.walletsRepo.SaveExtras(msg, tx.TxHash)
}
