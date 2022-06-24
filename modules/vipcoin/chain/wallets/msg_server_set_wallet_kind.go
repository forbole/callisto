package wallets

import (
	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSetKind allows to properly handle a MsgSetKind
func (m *Module) handleMsgSetKind(tx *juno.Tx, index int, msg *typeswallets.MsgSetWalletKind) error {
	if err := m.walletsRepo.SaveKinds(msg, tx.TxHash); err != nil {
		return err
	}

	wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.Address))
	switch {
	case err != nil:
		return err
	case len(wallets) != 1:
		return typeswallets.ErrInvalidAddressField
	}

	wallets[0].Kind = msg.Kind

	return m.walletsRepo.UpdateWallets(wallets...)
}
