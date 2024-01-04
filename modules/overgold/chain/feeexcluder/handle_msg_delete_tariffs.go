package feeexcluder

import (
	"errors"
	"strconv"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	juno "github.com/forbole/juno/v5/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgDeleteTariffs allows to properly handle a message
func (m *Module) handleMsgDeleteTariffs(tx *juno.Tx, _ int, msg *types.MsgDeleteTariffs) error {
	// 1) check if msg delete address already exists
	msgs, err := m.feeexcluderRepo.GetAllMsgDeleteTariffs(filter.NewFilter().SetArgument(db.FieldTxHash, tx.TxHash))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: "msg delete tariffs, creator: " + msg.Creator}
	}

	// 2) insert msg delete address to db
	if err = m.feeexcluderRepo.InsertToMsgDeleteTariffs(tx.TxHash, *msg); err != nil {
		return err
	}

	// 3) tariffs address by id
	tariffsID, err := strconv.ParseUint(msg.TariffID, 10, 64)
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return m.feeexcluderRepo.DeleteTariffs(nil, tariffsID)
}
