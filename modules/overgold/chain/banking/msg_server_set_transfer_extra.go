package banking

import (
	"git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSetTransferExtra allows to properly handle a handleMsgSetTransferExtra
func (m *Module) handleMsgSetTransferExtra(tx *juno.Tx, index int, msg *types.MsgSetTransferExtra) error {
	transfer, err := m.bankingRepo.GetBaseTransfers(filter.NewFilter().SetArgument(dbtypes.FieldID, msg.Id))
	switch {
	case err != nil:
		return err
	case len(transfer) != 1:
		return types.ErrInvalidAddressField
	}

	transfer[0].Extras = msg.Extras

	if err := m.bankingRepo.SaveBaseTransfers(transfer...); err != nil {
		return err
	}

	return m.bankingRepo.SaveMsgSetTransferExtra(msg, tx.TxHash)
}
