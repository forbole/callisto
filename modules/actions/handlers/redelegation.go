package handlers

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/modules/actions/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"
)

func RedelegationHandler(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msg("executing redelegations action")

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

	redelegationsList := make([]types.Redelegation, len(redelegations.RedelegationResponses))
	for index, del := range redelegations.RedelegationResponses {
		redelegationsList[index] = types.Redelegation{
			DelegatorAddress:    del.Redelegation.DelegatorAddress,
			ValidatorSrcAddress: del.Redelegation.ValidatorSrcAddress,
			ValidatorDstAddress: del.Redelegation.ValidatorDstAddress,
		}

		RedelegationEntriesList := make([]types.RedelegationEntry, len(del.Entries))
		for indexEntry, entry := range del.Entries {
			RedelegationEntriesList[indexEntry] = types.RedelegationEntry{
				CompletionTime: entry.RedelegationEntry.CompletionTime,
				Balance:        entry.Balance,
			}
		}

		redelegationsList[index].RedelegationEntries = RedelegationEntriesList
	}

	return types.RedelegationResponse{
		Redelegations: redelegationsList,
		Pagination:    redelegations.Pagination,
	}, nil
}
