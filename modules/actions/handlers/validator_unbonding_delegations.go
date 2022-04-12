package handlers

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/modules/actions/types"

	"github.com/rs/zerolog/log"
)

func ValidatorUnbondingDelegationsHandler(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msg("executing validator unbonding delegations action")

	// Get latest node height
	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get all unbonding delegations from the given validator opr address
	unbondingDelegations, err := ctx.Sources.StakingSource.GetUnbondingDelegationsFromValidator(
		height,
		payload.GetAddress(),
		payload.GetPagination(),
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting all unbonding delegations from validator %s: %s",
			payload.GetAddress(), err)
	}

	unbondingDelegationsList := make([]types.UnbondingDelegation, len(unbondingDelegations.UnbondingResponses))
	for index, del := range unbondingDelegations.UnbondingResponses {
		unbondingDelegationsList[index] = types.UnbondingDelegation{
			DelegatorAddress: del.DelegatorAddress,
			ValidatorAddress: del.ValidatorAddress,
			Entries:          del.Entries,
		}
	}

	return types.UnbondingDelegationResponse{
		UnbondingDelegations: unbondingDelegationsList,
		Pagination:           unbondingDelegations.Pagination,
	}, nil
}
