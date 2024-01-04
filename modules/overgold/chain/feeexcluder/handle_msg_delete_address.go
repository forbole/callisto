package feeexcluder

import (
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	juno "github.com/forbole/juno/v5/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgDeleteAddress allows to properly handle a message
func (m *Module) handleMsgDeleteAddress(tx *juno.Tx, _ int, msg *types.MsgDeleteAddress) error {
	// 1) check if msg delete address already exists
	msgs, err := m.feeexcluderRepo.GetAllMsgDeleteAddress(filter.NewFilter().SetArgument(db.FieldTxHash, tx.TxHash))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: "msg delete address, creator: " + msg.Creator}
	}

	// 2) insert msg delete address to db
	if err = m.feeexcluderRepo.InsertToMsgDeleteAddress(tx.TxHash, *msg); err != nil {
		return err
	}

	// 3) delete address by id
	return m.feeexcluderRepo.DeleteAddress(nil, msg.Id)
}
