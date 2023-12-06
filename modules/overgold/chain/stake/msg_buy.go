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

// handleMsgBuy allows to properly handle a stake buy message
func (m *Module) handleMsgBuy(tx *juno.Tx, _ int, msg *types.MsgBuyRequest) error {
	msgs, err := m.stakeRepo.GetAllMsgBuy(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldTxHash, tx.TxHash).
		SetArgument(db.FieldCreator, msg.Creator))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: fmt.Sprintf("msg buy, creator: %s, buy amount %s  ", msg.Creator, msg.Amount)}
	}

	return m.stakeRepo.InsertMsgBuy(tx.TxHash, types.MsgBuyRequest{
		Creator: msg.Creator,
		Amount:  msg.Amount,
	})
}
