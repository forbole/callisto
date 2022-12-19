package top_accounts

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	err := m.refreshBalances(msg, tx)
	if err != nil {
		return fmt.Errorf("error while refreshing account available balances: %s", err)
	}

	switch cosmosMsg := msg.(type) {
	// Handle x/staking delegations, redelegations, and unbondings
	case *stakingtypes.MsgDelegate:
		return m.handleMsgDelegate(tx.Height, cosmosMsg)

		// case *stakingtypes.MsgBeginRedelegate:
		// 	return m.stakingModule.HandleMsgBeginRedelegate(tx, index, cosmosMsg)

		// case *stakingtypes.MsgUndelegate:
		// 	return m.stakingModule.HandleMsgUndelegate(tx, index, cosmosMsg)

		// // Handle x/distribution delegator rewards
		// case *distritypes.MsgWithdrawDelegatorReward:
		// 	err := m.distrModule.RefreshDelegatorRewards([]string{cosmosMsg.DelegatorAddress}, tx.Height)
		// 	if err != nil {
		// 		return fmt.Errorf("error while refreshing delegator rewards from message: %s", err)
		// 	}

	}

	return nil
}

func (m *Module) refreshBalances(msg sdk.Msg, tx *juno.Tx) error {

	addresses, err := m.messageParser(m.cdc, msg)
	if err != nil {
		return fmt.Errorf("error while parsing account addresses of message type %s: %s", proto.MessageName(msg), err)
	}

	err = m.bankModule.UpdateBalances(utils.FilterNonAccountAddresses(addresses), tx.Height)
	if err != nil {
		return fmt.Errorf("error while updating account available balances: %s", err)
	}

	err = m.refreshTopAccountsSum(addresses)
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum while refreshing balance: %s", err)
	}

	return nil
}

func (m *Module) handleMsgDelegate(height int64, msg *stakingtypes.MsgDelegate) error {
	err := m.stakingModule.HandleMsgDelegate(height, msg)
	if err != nil {
		return fmt.Errorf("error while handling MsgDelegate: %s", err)
	}

	err = m.refreshTopAccountsSum([]string{msg.DelegatorAddress})
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum while handling MsgDelegate: %s", err)
	}

	return nil
}
