package staking

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc/codes"

	"github.com/forbole/bdjuno/v2/types"
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

// ConvertDelegationsResponses converts the given responses to BDJuno Delegation instances
func ConvertDelegationsResponses(height int64, responses []stakingtypes.DelegationResponse) []types.Delegation {
	var delegations = make([]types.Delegation, len(responses))
	for index, delegation := range responses {
		delegations[index] = convertDelegationResponse(height, delegation)
	}
	return delegations
}

// --------------------------------------------------------------------------------------------------------------------

// getValidatorDelegations returns all the delegations associated to the given validator at a given height
func (m *Module) getValidatorDelegations(height int64, validator string) ([]types.Delegation, error) {
	delegations, err := m.source.GetValidatorDelegations(height, validator)
	if err != nil {
		return nil, fmt.Errorf("error while getting validator delegations: %s", err)
	}

	return ConvertDelegationsResponses(height, delegations), nil
}

// getDelegatorDelegations returns the current delegations for the given delegator
func (m *Module) getDelegatorDelegations(height int64, delegator string) ([]types.Delegation, error) {
	responses, err := m.source.GetDelegatorDelegations(height, delegator)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator delegations: %s", err)
	}

	return ConvertDelegationsResponses(height, responses), nil
}

// --------------------------------------------------------------------------------------------------------------------

// RefreshValidatorDelegations refreshes the delegations for the validator with the given consensus address
func (m *Module) RefreshValidatorDelegations(height int64, valOperAddr string) error {
	responses, err := m.source.GetValidatorDelegations(height, valOperAddr)
	if err != nil {
		return fmt.Errorf("error while getting validator delegations: %s", err)
	}

	err = m.db.DeleteValidatorDelegations(valOperAddr)
	if err != nil {
		return fmt.Errorf("error while deleting validator delegations: %s", err)
	}

	err = m.db.SaveDelegations(ConvertDelegationsResponses(height, responses))
	if err != nil {
		return fmt.Errorf("error while storing validator delegations: %s", err)
	}

	return nil
}

// refreshDelegatorDelegations updates the delegations of the given delegator by querying them at the
// required height, and then stores them inside the database by replacing all existing ones.
func (m *Module) refreshDelegatorDelegations(height int64, delegator string) error {
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
