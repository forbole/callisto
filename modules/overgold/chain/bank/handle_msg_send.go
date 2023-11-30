package bank

import (
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgSend allows to properly handle a MsgSend
func (m *Module) handleMsgSend(tx *juno.Tx, index int, msg *bank.MsgSend) error {
	// 1) check if already exists (not found is ok)
	msgs, err := m.bankRepo.GetAllMsgSend(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(types.FieldTxHash, tx.TxHash).
		SetArgument(types.FieldFromAddress, msg.FromAddress).
		SetArgument(types.FieldToAddress, msg.ToAddress))
	if err != nil && !errors.Is(err, errs.NotFound{}) {
		return err
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: "msg send, hash: " + tx.TxHash}
	}

	// 2) insert to table
	return m.bankRepo.InsertMsgSend(tx.TxHash, bank.MsgSend{
		FromAddress: msg.FromAddress,
		ToAddress:   msg.ToAddress,
		Amount:      msg.Amount,
	})
}
