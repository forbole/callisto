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

// UpdateValidatorsUnbondingDelegations updates the redelegations for all the validators provided
func UpdateValidatorsRedelegations(
	height int64, validators []stakingtypes.Validator, client stakingtypes.QueryClient, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators redelegations")

	params, err := db.GetStakingParams()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var out = make(chan types.Redelegation)
	for _, val := range validators {
		wg.Add(1)
		go getRedelegations(val.OperatorAddress, params.BondName, height, client, out, &wg)
	}

	// We need to call wg.Wait inside another goroutine in order to solve the hanging bug that's described here:
	// https://dev.to/sophiedebenedetto/synchronizing-go-routines-with-channels-and-waitgroups-3ke2
	go func() {
		wg.Wait()
		close(out)
	}()

	var redelegations []types.Redelegation
	for del := range out {
		redelegations = append(redelegations, del)
	}

	return db.SaveRedelegations(redelegations)
}

func getRedelegations(
	validatorAddress string, bondDenom string, height int64, client stakingtypes.QueryClient,
	out chan<- types.Redelegation, wg *sync.WaitGroup,
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
			log.Error().Str("module", "staking").Err(err).
				Msg("error while getting validators redelegations")
			return
		}

		for _, delegation := range res.RedelegationResponses {
			for _, entry := range delegation.Entries {
				out <- types.NewRedelegation(
					delegation.Redelegation.DelegatorAddress,
					delegation.Redelegation.ValidatorSrcAddress,
					delegation.Redelegation.ValidatorDstAddress,
					sdk.NewCoin(bondDenom, entry.Balance),
					entry.RedelegationEntry.CompletionTime,
					entry.RedelegationEntry.CreationHeight,
				)
			}
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
	}
}
