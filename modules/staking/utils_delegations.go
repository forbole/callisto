package staking

import (
	"fmt"

	"google.golang.org/grpc/codes"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/forbole/bdjuno/types"
)

const (
	ErrDelegationNotFound = "rpc error: code = %s desc = rpc error: code = %s"
)

// convertDelegationResponse converts the given response to a BDJuno Delegation instance
func convertDelegationResponse(height int64, response stakingtypes.DelegationResponse) types.Delegation {
	return types.NewDelegation(
		response.Delegation.DelegatorAddress,
		response.Delegation.ValidatorAddress,
		response.Balance,
		height,
	)
}

// --------------------------------------------------------------------------------------------------------------------

// getDelegatorDelegations returns the current delegations for the given delegator
func (m *Module) getDelegatorDelegations(height int64, delegator string) ([]types.Delegation, error) {
	// Get the delegations
	responses, err := m.source.GetDelegatorDelegations(height, delegator)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator delegations: %s", err)
	}

	var delegations = make([]types.Delegation, len(responses))
	for index, delegation := range responses {
		delegations[index] = convertDelegationResponse(height, delegation)
	}

	return delegations, nil
}

// --------------------------------------------------------------------------------------------------------------------

// refreshDelegations updates the delegations of the given delegator by querying them at the
// required height, and then stores them inside the database by replacing all existing ones.
func (m *Module) refreshDelegations(height int64, delegator string) error {
	// Get current delegations
	delegations, err := m.getDelegatorDelegations(height, delegator)
	if err != nil {
		// Get the error code
		var code string
		_, scanErr := fmt.Sscanf(err.Error(), ErrDelegationNotFound, &code, &code)
		if scanErr != nil {
			return fmt.Errorf("error while scanning error: %s", scanErr)
		}

		// If delegations are not found there is no problem.
		// If it's a different error, we need to return it
		if code != codes.NotFound.String() {
			return fmt.Errorf("error while getting delegator delegations: %s", err)
		}
	}

	// Remove existing delegations
	err = m.db.DeleteDelegatorDelegations(delegator)
	if err != nil {
		return fmt.Errorf("error while deleting delegator delegations: %s", err)
	}

	// Save new delegations
	err = m.db.SaveDelegations(delegations)
	if err != nil {
		return fmt.Errorf("error while saving delegations: %s", err)
	}

	// Refresh the delegator rewards
	return m.distrModule.RefreshDelegatorRewards(height, delegator)
}
