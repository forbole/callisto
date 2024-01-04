package bank

import (
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgMultiSend allows to properly handle a MsgMultiSend
func (m *Module) handleMsgMultiSend(tx *juno.Tx, _ int, msg *bank.MsgMultiSend) error {
	// 1) check if already exists (not found is ok)
	msgs, err := m.bankRepo.GetAllMsgMultiSend(filter.NewFilter().SetArgument(types.FieldTxHash, tx.TxHash))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: "msg multi send, hash: " + tx.TxHash}
	}

	// 2) insert to table
	return m.bankRepo.InsertMsgMultiSend(tx.TxHash, bank.MsgMultiSend{
		Inputs:  msg.Inputs,
		Outputs: msg.Outputs,
	})
}
