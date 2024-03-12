package feeexcluder

import (
	"strconv"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgDeleteTariffs allows to properly handle a message
func (m *Module) handleMsgDeleteTariffs(tx *juno.Tx, _ int, msg *types.MsgDeleteTariffs) error {
	if err := m.feeexcluderRepo.InsertToMsgDeleteTariffs(tx.TxHash, *msg); err != nil {
		return err
	}

	tariffsID, err := strconv.ParseUint(msg.TariffID, 10, 64)
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return m.feeexcluderRepo.DeleteTariffs(nil, tariffsID)
}
