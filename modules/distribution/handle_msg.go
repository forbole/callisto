package distribution

import (
	"context"
	"fmt"

	bankutils "github.com/forbole/bdjuno/modules/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/desmos-labs/juno/client"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/distribution/utils"
)

// HandleMsg allows to handle the different utils related to the distribution module
func HandleMsg(
	tx *juno.Tx, msg sdk.Msg, distrClient distrtypes.QueryClient, bankClient banktypes.QueryClient, db *database.Db,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	if _, ok := msg.(*distrtypes.MsgFundCommunityPool); ok {
		return utils.UpdateCommunityPool(tx.Height, distrClient, db)
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

	header := client.GetHeightRequestHeader(tx.Height)
	res, err := distrClient.DelegatorWithdrawAddress(
		context.Background(),
		&distrtypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: delegatorAddr},
		header,
	)
	if err != nil {
		return fmt.Errorf("error while getting delegator withdraw address: %s", err)
	}

	var addresses = []string{delegatorAddr}
	if delegatorAddr != res.WithdrawAddress {
		// Only update the withdraw address if it's not the same as the delegation address
		addresses = append(addresses, res.WithdrawAddress)
	}

	return bankutils.RefreshBalances(tx.Height, addresses, bankClient, db)
}
