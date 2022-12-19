package top_accounts

import (
	"fmt"
	"time"

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

	// Refresh x/bank available account balances
	err := m.refreshBalances(msg, tx)
	if err != nil {
		return fmt.Errorf("error while refreshing account available balances: %s", err)
	}

	// Handle x/staking delegations, redelegations, and unbondings
	switch cosmosMsg := msg.(type) {
	case *stakingtypes.MsgDelegate:
		return m.handleMsgDelegate(tx.Height, cosmosMsg)

	case *stakingtypes.MsgBeginRedelegate:
		return m.handleMsgBeginRedelegate(tx, index, cosmosMsg)

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
	err := m.stakingModule.RefreshDelegations(height, msg.DelegatorAddress)
	if err != nil {
		return fmt.Errorf("error while refreshing delegations: %s", err)
	}

	err = m.refreshTopAccountsSum([]string{msg.DelegatorAddress})
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum while handling MsgDelegate: %s", err)
	}

	return nil
}

func (m *Module) handleMsgBeginRedelegate(
	tx *juno.Tx, index int, msg *stakingtypes.MsgBeginRedelegate) error {

	err := m.stakingModule.RefreshRedelegations(tx, index, msg.DelegatorAddress)
	if err != nil {
		return fmt.Errorf("error while getting redelegations while handling MsgBeginRedelegate: %s", err)
	}

	err = m.refreshTopAccountsSum([]string{msg.DelegatorAddress})
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum while handling MsgDelegate: %s", err)
	}

	event, err := tx.FindEventByType(index, stakingtypes.EventTypeRedelegate)
	if err != nil {
		return err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return err
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return err
	}

	// When the time expires, update the delegations & redelegations
	time.AfterFunc(time.Until(completionTime), m.refreshDelegations(tx.Height, msg.DelegatorAddress))
	time.AfterFunc(time.Until(completionTime), m.refreshRedelegations(tx, index, msg.DelegatorAddress))

	return nil
}
