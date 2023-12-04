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

// handleMsgSell allows to properly handle a stake sell message
func (m *Module) handleMsgSell(tx *juno.Tx, index int, msg *types.MsgSellRequest) error {
	msgs, err := m.stakeRepo.GetAllMsgSell(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldTxHash, tx.TxHash).
		SetArgument(db.FieldCreator, msg.Creator))
	if err != nil && !errors.Is(err, errs.NotFound{}) {
		return err
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: fmt.Sprintf("msg sell, creator: %s, sell amount %s  ", msg.Creator, msg.Amount)}
	}

	return m.stakeRepo.InsertMsgSell(tx.TxHash, types.MsgSellRequest{
		Creator: msg.Creator,
		Amount:  msg.Amount,
	})
}
