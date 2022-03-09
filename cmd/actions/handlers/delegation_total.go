package handlers

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func TotalDelegationAmountHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get all  delegations for given delegator address
	delegationList, err := ctx.Sources.StakingSource.GetDelegationsWithPagination(height, payload.GetAddress(), nil)
	if err != nil {
		// For stargate only, returns empty object without error if delegator delegations are not found on the chain
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

	return actionstypes.Balance{
		Coins: actionstypes.ConvertCoins(coinObject),
	}, nil
}
