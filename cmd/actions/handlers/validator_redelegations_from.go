package handlers

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func ValidatorRedelegationsFromHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	log.Debug().Str("action", "redelegation from validator").
		Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msgf("pagination query: %v", payload.GetPagination())

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get redelegations from a source validator address
	redelegations, err := ctx.Sources.StakingSource.GetRedelegations(height, &stakingtypes.QueryRedelegationsRequest{
		SrcValidatorAddr: payload.GetAddress(),
		Pagination:       payload.GetPagination(),
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting redelegations from validator: %s", err)
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
