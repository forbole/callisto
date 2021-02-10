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

// UpdateValidatorsUnbondingDelegations updates the unbonding delegations for all the validators provided
func UpdateValidatorsUnbondingDelegations(
	height int64, bondDenom string, validators []stakingtypes.Validator,
	client stakingtypes.QueryClient, db *database.BigDipperDb,
) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators unbonding delegations")

	var wg sync.WaitGroup
	for _, val := range validators {
		wg.Add(1)
		go getUnbondingDelegations(val.OperatorAddress, bondDenom, height, client, db, &wg)
	}
	wg.Wait()
}

// getUnbondingDelegations gets all the unbonding delegations referring to the validator having the
// given address at the given block height (having the given timestamp).
// All the unbonding delegations will be sent to the out channel, and wg.Done() will be called at the end.
func getUnbondingDelegations(
	validatorAddress string, bondDenom string, height int64,
	client stakingtypes.QueryClient, db *database.BigDipperDb, wg *sync.WaitGroup,
) {
	defer wg.Done()

	header := utils.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := client.ValidatorUnbondingDelegations(
			context.Background(),
			&stakingtypes.QueryValidatorUnbondingDelegationsRequest{
				ValidatorAddr: validatorAddress,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 unbonding delegations at a time
				},
			},
			header,
		)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).
				Int64("height", height).Str("validator", validatorAddress).
				Msg("error while getting validator delegations")
			return
		}

		var delegations []types.UnbondingDelegation
		for _, delegation := range res.UnbondingResponses {
			for _, entry := range delegation.Entries {
				delegations = append(delegations, types.NewUnbondingDelegation(
					delegation.DelegatorAddress,
					delegation.ValidatorAddress,
					sdk.NewCoin(bondDenom, entry.Balance),
					entry.CompletionTime,
					height,
				))
			}
		}

		err = db.SaveUnbondingDelegations(delegations)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).
				Int64("height", height).Str("validator", validatorAddress).
				Msg("error while saving validator delegations")
			return
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
	}
}
