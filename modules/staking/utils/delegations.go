package utils

import (
	"context"
	"fmt"
	"sync"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"google.golang.org/grpc/codes"

	distrutils "github.com/forbole/bdjuno/modules/distribution/utils"

	"github.com/desmos-labs/juno/client"

	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"
)

const (
	ErrDelegationNotFound = "rpc error: code = %s desc = rpc error: code = %s"
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
	height int64, validators []stakingtypes.Validator, client stakingtypes.QueryClient, db *database.Db,
) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators delegations")

	var wg sync.WaitGroup
	for _, val := range validators {
		wg.Add(1)
		go getDelegationsFromGrpc(val.OperatorAddress, height, client, db, &wg)
	}
	wg.Wait()
}

// getDelegationsFromGrpc gets the list of all the delegations that the validator having the given address has
// at the given block height (having the given timestamp).
// All the delegations will be sent to the out channel, and wg.Done() will be called at the end.
func getDelegationsFromGrpc(
	validatorAddress string, height int64, stakingClient stakingtypes.QueryClient, db *database.Db, wg *sync.WaitGroup,
) {
	defer wg.Done()

	header := client.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := stakingClient.ValidatorDelegations(
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
			log.Error().Str("module", "staking").Err(err).
				Int64("height", height).
				Str("validator", validatorAddress).
				Msg("error while saving validator delegations")
			return
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
	}
}

// --------------------------------------------------------------------------------------------------------------------

// GetDelegatorDelegations returns the current delegations for the given delegator
func GetDelegatorDelegations(height int64, delegator string, client stakingtypes.QueryClient) ([]types.Delegation, error) {
	// Get the delegations
	res, err := client.DelegatorDelegations(
		context.Background(),
		&stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: delegator,
		},
	)
	if err != nil {
		return nil, err
	}

	var delegations = make([]types.Delegation, len(res.DelegationResponses))
	for index, delegation := range res.DelegationResponses {
		delegations[index] = ConvertDelegationResponse(height, delegation)
	}

	return delegations, nil
}

// --------------------------------------------------------------------------------------------------------------------

// RefreshDelegations updates the delegations of the given delegator by querying them at the
// required height, and then stores them inside the database by replacing all existing ones.
func RefreshDelegations(
	height int64, delegator string,
	stakingClient stakingtypes.QueryClient, distrClient distrtypes.QueryClient,
	db *database.Db,
) error {
	// Get current delegations
	delegations, err := GetDelegatorDelegations(height, delegator, stakingClient)
	if err != nil {
		// Get the error code
		var code string
		_, scanErr := fmt.Sscanf(err.Error(), ErrDelegationNotFound, &code, &code)
		if scanErr != nil {
			return err
		}

		// If delegations are not found there is no problem.
		// If it's a different error, we need to return it
		if code != codes.NotFound.String() {
			return fmt.Errorf("error while getting delegator delegations: %s", err)
		}
	}

	// Remove existing delegations
	err = db.DeleteDelegatorDelegations(delegator)
	if err != nil {
		return fmt.Errorf("error while deleting delegator delegations: %s", err)
	}

	// Save new delegations
	err = db.SaveDelegations(delegations)
	if err != nil {
		return fmt.Errorf("error while saving delegations: %s", err)
	}

	// Refresh the delegator rewards
	return distrutils.RefreshDelegatorRewards(height, delegator, distrClient, db)
}
