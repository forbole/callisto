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
) error {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators delegations")

	var wg sync.WaitGroup
	var out = make(chan types.Delegation)
	for _, val := range validators {
		wg.Add(1)
		go getDelegations(val.OperatorAddress, height, client, out, &wg)
	}

	// We need to call wg.Wait inside another goroutine in order to solve the hanging bug that's described here:
	// https://dev.to/sophiedebenedetto/synchronizing-go-routines-with-channels-and-waitgroups-3ke2
	go func() {
		wg.Wait()
		close(out)
	}()

	var delegations []types.Delegation
	for delegation := range out {
		delegations = append(delegations, delegation)
	}

	return db.SaveDelegations(delegations)
}

// getDelegations gets the list of all the delegations that the validator having the given address has
// at the given block height (having the given timestamp).
// All the delegations will be sent to the out channel, and wg.Done() will be called at the end.
func getDelegations(
	validatorAddress string, height int64, client stakingtypes.QueryClient,
	out chan<- types.Delegation, wg *sync.WaitGroup,
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

		for _, delegation := range res.DelegationResponses {
			out <- types.NewDelegation(
				delegation.Delegation.DelegatorAddress,
				delegation.Delegation.ValidatorAddress,
				delegation.Balance,
				delegation.Delegation.Shares.String(),
				height,
			)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
	}
}
