package feeexcluder

import (
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	juno "github.com/forbole/juno/v5/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgCreateTariffs allows to properly handle a message
func (m *Module) handleMsgCreateTariffs(tx *juno.Tx, _ int, msg *types.MsgCreateTariffs) error {
	// 1) check if msg create tariffs already exists
	msgs, err := m.feeexcluderRepo.GetAllMsgCreateTariffs(filter.NewFilter().SetArgument(db.FieldTxHash, tx.TxHash))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: "msg create tariffs, creator: " + msg.Creator}
	}

	// 2) insert msg create tariffs to db (including inserting into tariffs table)
	return m.feeexcluderRepo.InsertToMsgCreateTariffs(tx.TxHash, *msg)
}
