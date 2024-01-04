package feeexcluder

import (
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	juno "github.com/forbole/juno/v5/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgCreateAddress allows to properly handle a message
func (m *Module) handleMsgCreateAddress(tx *juno.Tx, _ int, msg *types.MsgCreateAddress) error {
	// 1) check if msg create address already exists
	msgs, err := m.feeexcluderRepo.GetAllMsgCreateAddress(filter.NewFilter().SetArgument(db.FieldTxHash, tx.TxHash))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: "msg create address, address: " + msg.Address}
	}

	// 2) insert msg create address to db (including inserting into address table)
	return m.feeexcluderRepo.InsertToMsgCreateAddress(tx.TxHash, *msg)
}
