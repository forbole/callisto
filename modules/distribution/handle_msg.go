package distribution

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	if _, ok := msg.(*distrtypes.MsgFundCommunityPool); ok {
		return m.updateCommunityPool(tx.Height)
	}

	var delegatorAddr string
	switch cosmosMsg := msg.(type) {
	case *distrtypes.MsgWithdrawValidatorCommission:
		delegatorAddr = cosmosMsg.GetSigners()[0].String()
	case *distrtypes.MsgWithdrawDelegatorReward:
		delegatorAddr = cosmosMsg.DelegatorAddress
	default:
		return nil
	}

	withdrawAddr, err := m.source.DelegatorWithdrawAddress(delegatorAddr, tx.Height)
	if err != nil {
		return fmt.Errorf("error while getting delegator withdraw address: %s", err)
	}

	var addresses = []string{delegatorAddr}
	if delegatorAddr != withdrawAddr {
		// Only update the withdraw address if it's not the same as the delegation address
		addresses = append(addresses, withdrawAddr)
	}

	return m.bankModule.RefreshBalances(tx.Height, addresses)
}
