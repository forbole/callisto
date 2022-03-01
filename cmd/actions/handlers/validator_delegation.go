package handlers

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func ValidatorDelegation(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	// Get latest node height
	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get validator's total delegations
	res, err := ctx.Sources.StakingSource.GetValidatorDelegationsWithPagination(height, payload.GetAddress(), payload.GetPagination())
	if err != nil {
		return nil, fmt.Errorf("error while getting validator delegations: %s", err)
	}

	delegations := make([]actionstypes.Delegation, len(res.DelegationResponses))
	for index, del := range res.DelegationResponses {
		delegations[index] = actionstypes.Delegation{
			DelegatorAddress: del.Delegation.DelegatorAddress,
			ValidatorAddress: del.Delegation.ValidatorAddress,
			Coins:            actionstypes.ConvertCoins([]sdk.Coin{del.Balance}),
		}
	}

	return actionstypes.DelegationResponse{
		Delegations: delegations,
		Pagination:  res.Pagination,
	}, nil
}
