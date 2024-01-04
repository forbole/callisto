package handlers

import (
	"fmt"
	"strings"

	"github.com/forbole/bdjuno/v4/modules/actions/types"

	"google.golang.org/grpc/codes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
)

func DelegationHandler(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Str("action", "delegations").
		Str("address", payload.GetAddress()).
		Msg("executing delegations action")

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get delegator's total rewards
	res, err := ctx.Sources.StakingSource.GetDelegationsWithPagination(height, payload.GetAddress(), payload.GetPagination())
	if err != nil {
		// For stargate only, returns without throwing error if delegator delegations are not found on the chain
		if strings.Contains(err.Error(), codes.NotFound.String()) {
			return err, nil
		}
		return err, fmt.Errorf("error while getting delegator delegations: %s", err)
	}

	delegations := make([]types.Delegation, len(res.DelegationResponses))
	for index, del := range res.DelegationResponses {
		delegations[index] = types.Delegation{
			DelegatorAddress: del.Delegation.DelegatorAddress,
			ValidatorAddress: del.Delegation.ValidatorAddress,
			Coins:            types.ConvertCoins([]sdk.Coin{del.Balance}),
		}
	}

	return types.DelegationResponse{
		Delegations: delegations,
		Pagination:  res.Pagination,
	}, nil
}
