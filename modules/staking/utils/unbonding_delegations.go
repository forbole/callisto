package utils

import (
	"context"
	"sync"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"
)

// ConvertUnbondingResponse converts the given UnbondingDelegation response into a slice of BDJuno UnbondingDelegation
func ConvertUnbondingResponse(
	height int64, bondDenom string, response stakingtypes.UnbondingDelegation,
) []types.UnbondingDelegation {
	var delegations []types.UnbondingDelegation
	for _, entry := range response.Entries {
		delegations = append(delegations, types.NewUnbondingDelegation(
			response.DelegatorAddress,
			response.ValidatorAddress,
			sdk.NewCoin(bondDenom, entry.Balance),
			entry.CompletionTime,
			height,
		))
	}
	return delegations
}

// --------------------------------------------------------------------------------------------------------------------

// UpdateValidatorsUnbondingDelegations updates the unbonding delegations for all the validators provided
func UpdateValidatorsUnbondingDelegations(
	height int64, bondDenom string, validators []stakingtypes.Validator,
	client stakingtypes.QueryClient, db *database.Db,
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
	stakingClient stakingtypes.QueryClient, db *database.Db, wg *sync.WaitGroup,
) {
	defer wg.Done()

	header := client.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := stakingClient.ValidatorUnbondingDelegations(
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
			delegations = append(delegations, ConvertUnbondingResponse(height, bondDenom, delegation)...)
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

// --------------------------------------------------------------------------------------------------------------------

// DeleteUnbondingDelegation returns a function that when called deletes the given delegation from the database
func DeleteUnbondingDelegation(delegation types.UnbondingDelegation, db *database.Db) func() {
	return func() {
		err := db.DeleteUnbondingDelegation(delegation)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).
				Str("operation", "delete unbonding delegation").Msg("error while deleting unbonding delegation")
		}
	}
}
