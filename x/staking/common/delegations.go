package common

import (
	"context"
	"sync"

	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	"github.com/forbole/bdjuno/x/utils"
)

// UpdateValidatorsDelegations updates the delegations for all the given validators at the provided height
func UpdateValidatorsDelegations(
	height int64, validators []stakingtypes.Validator, client stakingtypes.QueryClient, db *database.BigDipperDb,
) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators delegations")

	var wg sync.WaitGroup
	for _, val := range validators {
		wg.Add(1)
		go getDelegations(val.OperatorAddress, height, client, db, &wg)
	}
	wg.Wait()
}

// getDelegations gets the list of all the delegations that the validator having the given address has
// at the given block height (having the given timestamp).
// All the delegations will be sent to the out channel, and wg.Done() will be called at the end.
func getDelegations(
	validatorAddress string, height int64,
	client stakingtypes.QueryClient, db *database.BigDipperDb, wg *sync.WaitGroup,
) {
	defer wg.Done()

	header := utils.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := client.ValidatorDelegations(
			context.Background(),
			&stakingtypes.QueryValidatorDelegationsRequest{
				ValidatorAddr: validatorAddress,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 delegations at time
				},
			},
			header,
		)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Str("validator", validatorAddress).Msg("error while getting validator delegations")
			return
		}

		var delegations = make([]types.Delegation, len(res.DelegationResponses))
		for index, delegation := range res.DelegationResponses {
			delegations[index] = types.NewDelegation(
				delegation.Delegation.DelegatorAddress,
				delegation.Delegation.ValidatorAddress,
				delegation.Balance,
				delegation.Delegation.Shares.String(),
				height,
			)
		}
		err = db.SaveDelegations(delegations)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Str("validator", validatorAddress).Msg("error while saving validator delegations")
			return
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
	}
}
