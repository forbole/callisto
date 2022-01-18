package staking

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

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
