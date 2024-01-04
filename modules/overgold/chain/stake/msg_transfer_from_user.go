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

// handleMsgTransferFromUser allows to properly handle a transfer from user message.
func (m *Module) handleMsgTransferFromUser(tx *juno.Tx, _ int, msg *types.MsgTransferFromUser) error {
	msgs, err := m.stakeRepo.GetAllMsgTransferFromUser(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldTxHash, tx.TxHash).
		SetArgument(db.FieldCreator, msg.Creator))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: fmt.Sprintf("msg transer from user, creator: %s, address: %s, amount %s",
			msg.Creator, msg.Address, msg.Amount)}
	}

	return m.stakeRepo.InsertMsgTransferFromUser(tx.TxHash, types.MsgTransferFromUser{
		Creator: msg.Creator,
		Amount:  msg.Amount,
		Address: msg.Address,
	})
}
