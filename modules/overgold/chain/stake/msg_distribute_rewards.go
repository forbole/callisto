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

// handleMsgDistributeRewards allows to properly handle a stake distribute rewards message
func (m *Module) handleMsgDistributeRewards(tx *juno.Tx, _ int, msg *types.MsgDistributeRewards) error {
	msgs, err := m.stakeRepo.GetAllMsgDistributeRewards(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldTxHash, tx.TxHash).
		SetArgument(db.FieldCreator, msg.Creator))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: fmt.Sprintf("msg distribute rewards, creator: %s", msg.Creator)}
	}

	return m.stakeRepo.InsertMsgDistributeRewards(tx.TxHash, types.MsgDistributeRewards{
		Creator: msg.Creator,
	})
}
