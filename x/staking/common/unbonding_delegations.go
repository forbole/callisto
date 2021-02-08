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

// UpdateValidatorUnbondingDelegations updates the unbonding delegations for all the validators provided
func UpdateValidatorUnbondingDelegations(
	height int64, validators []stakingtypes.Validator, client stakingtypes.QueryClient, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "staking").Int64("height", height).Msg("updating validators unbonding delegations")

	params, err := db.GetStakingParams()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var out = make(chan types.UnbondingDelegation)
	for _, val := range validators {
		wg.Add(1)
		go getUnbondingDelegations(val.OperatorAddress, params.BondName, height, client, out, &wg)
	}

	// We need to call wg.Wait inside another goroutine in order to solve the hanging bug that's described here:
	// https://dev.to/sophiedebenedetto/synchronizing-go-routines-with-channels-and-waitgroups-3ke2
	go func() {
		wg.Wait()
		close(out)
	}()

	var delegations []types.UnbondingDelegation
	for del := range out {
		delegations = append(delegations, del)
	}

	return db.SaveUnbondingDelegations(delegations)
}

// getUnbondingDelegations returns the list of all the unbonding delegations that the validator having the
// given address has at the given block height (having the given timestamp).
func getUnbondingDelegations(
	validatorAddress string, bondDenom string, height int64, client stakingtypes.QueryClient,
	out chan<- types.UnbondingDelegation, wg *sync.WaitGroup,
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

		for _, delegation := range res.UnbondingResponses {
			for _, entry := range delegation.Entries {
				out <- types.NewUnbondingDelegation(
					delegation.DelegatorAddress,
					delegation.ValidatorAddress,
					sdk.NewCoin(bondDenom, entry.Balance),
					entry.CompletionTime,
					height,
				)
			}
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
	}
}
