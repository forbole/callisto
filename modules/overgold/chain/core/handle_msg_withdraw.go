package core

import (
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	juno "github.com/forbole/juno/v5/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgWithdraw allows to properly handle a MsgWithdraw
func (m *Module) handleMsgWithdraw(tx *juno.Tx, index int, msg *types.MsgWithdraw) error {
	msgs, err := m.coreRepo.GetAllMsgWithdraw(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldTxHash, tx.TxHash).
		SetArgument(db.FieldCreator, msg.Creator))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: "msg withdraw, address: " + msg.Address}
	}

	return m.coreRepo.InsertMsgWithdraw(tx.TxHash, types.MsgWithdraw{
		Creator: msg.Creator,
		Amount:  msg.Amount,
		Denom:   msg.Denom,
		Address: msg.Address,
	})
}
