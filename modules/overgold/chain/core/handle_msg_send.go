package core

import (
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	juno "github.com/forbole/juno/v5/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgSend allows to properly handle a MsgSend
func (m *Module) handleMsgSend(tx *juno.Tx, index int, msg *types.MsgSend) error {
	msgs, err := m.coreRepo.GetAllMsgSend(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldTxHash, tx.TxHash).
		SetArgument(db.FieldCreator, msg.Creator))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: "msg send, address from: " + msg.From}
	}

	return m.coreRepo.InsertMsgSend(tx.TxHash, types.MsgSend{
		Creator: msg.Creator,
		From:    msg.From,
		To:      msg.To,
		Amount:  msg.Amount,
		Denom:   msg.Denom,
	})
}
