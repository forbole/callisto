package handlers

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func RedelegationHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get delegator's redelegations
	redelegations, err := ctx.Sources.StakingSource.GetRedelegations(height, &stakingtypes.QueryRedelegationsRequest{
		DelegatorAddr: payload.GetAddress(),
		Pagination:    payload.GetPagination(),
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator redelegations: %s", err)
	}

	redelegationsList := make([]actionstypes.Redelegation, len(redelegations.RedelegationResponses))
	for index, del := range redelegations.RedelegationResponses {
		redelegationsList[index] = actionstypes.Redelegation{
			DelegatorAddress:    del.Redelegation.DelegatorAddress,
			ValidatorSrcAddress: del.Redelegation.ValidatorSrcAddress,
			ValidatorDstAddress: del.Redelegation.ValidatorDstAddress,
		}

		RedelegationEntriesList := make([]actionstypes.RedelegationEntry, len(del.Entries))
		for indexEntry, entry := range del.Entries {
			RedelegationEntriesList[indexEntry] = actionstypes.RedelegationEntry{
				CompletionTime: entry.RedelegationEntry.CompletionTime,
				Balance:        entry.Balance,
			}
		}

		redelegationsList[index].RedelegationEntries = RedelegationEntriesList
	}

	return actionstypes.RedelegationResponse{
		Redelegations: redelegationsList,
		Pagination:    redelegations.Pagination,
	}, nil
}
