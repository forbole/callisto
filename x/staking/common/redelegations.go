package common

import (
	"context"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	"github.com/forbole/bdjuno/x/utils"
)

// UpdateValidatorsRedelegations updates the redelegations for all the validators provided
func UpdateValidatorsRedelegations(
	height int64, bondDenom string, validators []stakingtypes.Validator,
	client stakingtypes.QueryClient, db *database.BigDipperDb,
) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators redelegations")

	var wg sync.WaitGroup
	for _, val := range validators {
		wg.Add(1)
		go getRedelegations(val.OperatorAddress, bondDenom, height, client, db, &wg)
	}
	wg.Wait()
}

func getRedelegations(
	validatorAddress string, bondDenom string, height int64,
	client stakingtypes.QueryClient, db *database.BigDipperDb, wg *sync.WaitGroup,
) {
	defer wg.Done()

	header := utils.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := client.Redelegations(
			context.Background(),
			&stakingtypes.QueryRedelegationsRequest{
				SrcValidatorAddr: validatorAddress,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 unbonding delegations at a time
				},
			},
			header,
		)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Msg("error while getting validators redelegations")
			return
		}

		var delegations []types.Redelegation
		for _, delegation := range res.RedelegationResponses {
			for _, entry := range delegation.Entries {
				delegations = append(delegations, types.NewRedelegation(
					delegation.Redelegation.DelegatorAddress,
					delegation.Redelegation.ValidatorSrcAddress,
					delegation.Redelegation.ValidatorDstAddress,
					sdk.NewCoin(bondDenom, entry.Balance),
					entry.RedelegationEntry.CompletionTime,
					entry.RedelegationEntry.CreationHeight,
				))
			}
		}
		err = db.SaveRedelegations(delegations)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Msg("error while saving validators redelegations")
			return
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
	}
}
