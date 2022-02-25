package handlers

import (
	"fmt"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func ValidatorUnbondingDelegationsHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
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
