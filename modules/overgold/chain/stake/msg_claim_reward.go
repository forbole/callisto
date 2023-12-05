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

// handleMsgClaimReward allows to properly handle a stake claim reward message
func (m *Module) handleMsgClaimReward(tx *juno.Tx, _ int, msg *types.MsgClaimReward) error {
	msgs, err := m.stakeRepo.GetAllMsgDistributeRewards(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldTxHash, tx.TxHash).
		SetArgument(db.FieldCreator, msg.Creator))
	if err != nil && !errors.Is(err, errs.NotFound{}) {
		return err
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: fmt.Sprintf("msg claim reward, creator: %s", msg.Creator)}
	}

	return m.stakeRepo.InsertMsgClaimReward(tx.TxHash, types.MsgClaimReward{
		Creator: msg.Creator,
	})
}
