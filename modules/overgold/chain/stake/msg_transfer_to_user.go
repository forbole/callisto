package stake

import (
	"errors"
	"fmt"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	juno "github.com/forbole/juno/v5/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgTransferToUser allows to properly handle a transfer to user message.
func (m *Module) handleMsgTransferToUser(tx *juno.Tx, _ int, msg *types.MsgTransferToUser) error {
	msgs, err := m.stakeRepo.GetAllMsgTransferToUser(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldTxHash, tx.TxHash).
		SetArgument(db.FieldCreator, msg.Creator))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: fmt.Sprintf("msg transer to user, creator: %s, address: %s, amount %s",
			msg.Creator, msg.Address, msg.Amount)}
	}

	return m.stakeRepo.InsertMsgTransferToUser(tx.TxHash, types.MsgTransferToUser{
		Creator: msg.Creator,
		Amount:  msg.Amount,
		Address: msg.Address,
	})
}
