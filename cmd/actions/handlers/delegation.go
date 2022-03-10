package handlers

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func DelegationHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
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
