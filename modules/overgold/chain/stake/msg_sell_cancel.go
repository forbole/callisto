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

// handleMsgSellCancel allows to properly handle a stake sell cancel message
func (m *Module) handleMsgSellCancel(tx *juno.Tx, _ int, msg *types.MsgMsgCancelSell) error {
	msgs, err := m.stakeRepo.GetAllMsgSellCancel(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldTxHash, tx.TxHash).
		SetArgument(db.FieldCreator, msg.Creator))
	if err != nil && !errors.Is(err, errs.NotFound{}) {
		return err
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: fmt.Sprintf("msg sell cancel, creator: %s, sell cancel amount %s  ",
			msg.Creator, msg.Amount.Amount.String())}
	}

	return m.stakeRepo.InsertMsgSellCancel(tx.TxHash, types.MsgMsgCancelSell{
		Creator: msg.Creator,
		Amount:  msg.Amount,
	})
}
