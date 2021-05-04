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

// ConvertDelegationResponse converts the given response to a BDJuno Delegation instance
func ConvertDelegationResponse(height int64, response stakingtypes.DelegationResponse) types.Delegation {
	return types.NewDelegation(
		response.Delegation.DelegatorAddress,
		response.Delegation.ValidatorAddress,
		response.Balance,
		height,
	)
}

// --------------------------------------------------------------------------------------------------------------------

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
			delegations[index] = ConvertDelegationResponse(height, delegation)
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

// --------------------------------------------------------------------------------------------------------------------

// UpdateDelegations updates the current delegations for the given delegator by removing all the existing ones and
// getting the new ones from the client
func UpdateDelegations(height int64, delegator string, client stakingtypes.QueryClient, db *database.BigDipperDb) error {
	// Remove existing delegations
	err := db.DeleteDelegatorDelegations(delegator)
	if err != nil {
		return err
	}

	// Get the delegations
	res, err := client.DelegatorDelegations(
		context.Background(),
		&stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: delegator,
		},
	)
	if err != nil {
		return err
	}

	var delegations = make([]types.Delegation, len(res.DelegationResponses))
	for index, delegation := range res.DelegationResponses {
		delegations[index] = ConvertDelegationResponse(height, delegation)
	}

	return db.SaveDelegations(delegations)
}

// RefreshDelegations returns a function that when called updates the delegations of the provided delegator.
// In order to properly update the data, it removes all the existing delegations and stores new ones querying the gRPC
func RefreshDelegations(height int64, delegator string, client stakingtypes.QueryClient, db *database.BigDipperDb) func() {
	return func() {
		err := UpdateDelegations(height, delegator, client, db)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).
				Str("operation", "refresh delegations").Msg("error while refreshing delegations")
		}
	}
}
