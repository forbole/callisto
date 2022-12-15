package top_accounts

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/bdjuno/v3/modules/utils"
	"github.com/forbole/bdjuno/v3/types"
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
	balances, err := m.bankModule.UpdateBalances(utils.FilterNonAccountAddresses(addresses), tx.Height)
	if err != nil {
		return fmt.Errorf("error while updating account available balances: %s", err)
	}

	// Store native token balances to top_accounts table
	err = m.saveTopAccountsAvailable(balances)
	if err != nil {
		return fmt.Errorf("error while saving available balances to top_accounts table: %s", err)
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

func (m *Module) saveTopAccountsAvailable(accountBalances []types.AccountBalance) error {
	balances := make([]types.NativeTokenBalance, len(accountBalances))
	addresses := make([]string, len(accountBalances))
	for index, bal := range accountBalances {
		balances[index] = types.NewNativeTokenBalance(
			bal.Address,
			bal.Balance.AmountOf("uqck"),
		)
		addresses[index] = bal.Address
	}

	err := m.db.SaveTopAccountsBalance("available", balances)
	if err != nil {
		return fmt.Errorf("error while saving top accounts available balances: %s", err)
	}

	return m.refreshTopAccountsSum(addresses)
}
