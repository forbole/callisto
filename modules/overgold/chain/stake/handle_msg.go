package stake

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch stakeMsg := msg.(type) {
	case *types.MsgSellRequest:
		return m.handleMsgSell(tx, index, stakeMsg)
	case *types.MsgMsgCancelSell:
		return m.handleMsgSellCancel(tx, index, stakeMsg)
	case *types.MsgBuyRequest:
		return m.handleMsgBuy(tx, index, stakeMsg)
	case *types.MsgDistributeRewards:
		return m.handleMsgDistributeRewards(tx, index, stakeMsg)
	case *types.MsgClaimReward:
		return m.handleMsgClaimReward(tx, index, stakeMsg)
	default:
		return nil
	}
}
