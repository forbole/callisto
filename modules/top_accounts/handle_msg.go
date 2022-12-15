package top_accounts

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/bdjuno/v3/modules/utils"
	juno "github.com/forbole/juno/v3/types"
	"github.com/gogo/protobuf/proto"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	// Refresh account available balances
	addresses, err := m.messageParser(m.cdc, msg)
	if err != nil {
		return fmt.Errorf("error while parsing account addresses of message type %s: %s", proto.MessageName(msg), err)
	}
	err = m.bankModule.UpdateBalances(utils.FilterNonAccountAddresses(addresses), tx.Height)
	if err != nil {
		return fmt.Errorf("error while updating account available balances: %s", err)
	}

	switch cosmosMsg := msg.(type) {

	// Handle x/staking delegations, redelegations, and unbondings
	case *stakingtypes.MsgDelegate:
		return m.stakingModule.HandleMsgDelegate(tx.Height, cosmosMsg)

	case *stakingtypes.MsgBeginRedelegate:
		return m.stakingModule.HandleMsgBeginRedelegate(tx, index, cosmosMsg)

	case *stakingtypes.MsgUndelegate:
		return m.stakingModule.HandleMsgUndelegate(tx, index, cosmosMsg)

	// Handle x/distribution delegator rewards
	case *distritypes.MsgWithdrawDelegatorReward:
		err := m.distrModule.RefreshDelegatorRewards([]string{cosmosMsg.DelegatorAddress}, tx.Height)
		if err != nil {
			return fmt.Errorf("error while refreshing delegator rewards from message: %s", err)
		}

	}

	return nil
}
