package handlers

import (
	"fmt"
	"strings"

	"github.com/forbole/callisto/v4/modules/actions/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
)

func TotalDelegationAmountHandler(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msg("executing total delegation amount action")

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get all  delegations for given delegator address
	delegationList, err := ctx.Sources.StakingSource.GetDelegationsWithPagination(height, payload.GetAddress(), nil)
	if err != nil {
		// For stargate only, returns without throwing error if delegator delegations are not found on the chain
		if strings.Contains(err.Error(), codes.NotFound.String()) {
			return err, nil
		}
		return err, fmt.Errorf("error while getting delegator delegations: %s", err)
	}

	var coinObject sdk.Coins

	// Add up total value of delegations
	for _, eachDelegation := range delegationList.DelegationResponses {
		for index, eachCoin := range coinObject {
			if eachCoin.Denom == eachDelegation.Balance.Denom {
				coinObject[index].Amount = coinObject[index].Amount.Add(eachDelegation.Balance.Amount)
			}
			if eachCoin.Denom != eachDelegation.Balance.Denom {
				coinObject = append(coinObject, sdk.NewCoin(eachDelegation.Balance.Denom, eachDelegation.Balance.Amount))
			}
		}
		if coinObject == nil {
			coinObject = append(coinObject, sdk.NewCoin(eachDelegation.Balance.Denom, eachDelegation.Balance.Amount))
		}
	}

	return types.Balance{
		Coins: types.ConvertCoins(coinObject),
	}, nil
}
