package handlers

import (
	"fmt"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/rs/zerolog/log"
)

func UnbondingDelegationsHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	log.Debug().Str("action", "unbonding delegations").
		Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msgf("pagination query: %v", payload.GetPagination())

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get all unbonding delegations for given delegator address
	unbondingDelegations, err := ctx.Sources.StakingSource.GetUnbondingDelegations(height, payload.GetAddress(), payload.GetPagination())
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator unbonding delegations: %s", err)
	}

	unbondingDelegationsList := make([]actionstypes.UnbondingDelegation, len(unbondingDelegations.UnbondingResponses))
	for index, del := range unbondingDelegations.UnbondingResponses {
		unbondingDelegationsList[index] = actionstypes.UnbondingDelegation{
			DelegatorAddress: del.DelegatorAddress,
			ValidatorAddress: del.ValidatorAddress,
			Entries:          del.Entries,
		}
	}

	return actionstypes.UnbondingDelegationResponse{
		UnbondingDelegations: unbondingDelegationsList,
		Pagination:           unbondingDelegations.Pagination,
	}, nil
}
