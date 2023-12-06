package core

import (
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	juno "github.com/forbole/juno/v5/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgIssue allows to properly handle a MsgIssue
func (m *Module) handleMsgIssue(tx *juno.Tx, index int, msg *types.MsgIssue) error {
	msgs, err := m.coreRepo.GetAllMsgIssue(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldTxHash, tx.TxHash).
		SetArgument(db.FieldCreator, msg.Creator))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: "msg issue, address: " + msg.Address}
	}

	return m.coreRepo.InsertMsgIssue(tx.TxHash, types.MsgIssue{
		Creator: msg.Creator,
		Amount:  msg.Amount,
		Denom:   msg.Denom,
		Address: msg.Address,
	})
}
