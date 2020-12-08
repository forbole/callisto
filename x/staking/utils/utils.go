package utils

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/client"
	"github.com/forbole/bdjuno/x/staking/types"
)

// GetDelegations returns the list of all the delegations that the validator having the given address has
// at the given block height (having the given timestamp)
func GetDelegations(
	validatorAddress sdk.ValAddress, height int64, timestamp time.Time, cp *client.Proxy,
) ([]types.Delegation, error) {
	var responses []staking.DelegationResponse
	endpoint := fmt.Sprintf("/staking/validators/%s/delegations?height=%d", validatorAddress.String(), height)
	if _, err := cp.QueryLCDWithHeight(endpoint, &responses); err != nil {
		return nil, err
	}

	delegations := make([]types.Delegation, len(responses))
	for index, delegation := range responses {
		delegations[index] = types.NewDelegation(
			delegation.GetDelegatorAddr(),
			delegation.GetValidatorAddr(),
			delegation.Balance,
			delegation.Shares.String(),
			height,
			timestamp,
		)
	}

	return delegations, nil
}

// GetUnbondingDelegations returns the list of all the unbonding delegations that the validator having the
// given address has at the given block height (having the given timestamp).
func GetUnbondingDelegations(
	validatorAddress sdk.ValAddress, bondDenom string, height int64, timestamp time.Time, cp *client.Proxy,
) ([]types.UnbondingDelegation, error) {
	var responses []staking.UnbondingDelegation
	endpoint := fmt.Sprintf("/staking/validators/%s/unbonding_delegations?height=%d", validatorAddress.String(), height)
	if _, err := cp.QueryLCDWithHeight(endpoint, &responses); err != nil {
		return nil, err
	}

	var unbondingDelegations []types.UnbondingDelegation
	for _, delegation := range responses {
		for _, entry := range delegation.Entries {
			unbondingDelegations = append(unbondingDelegations, types.NewUnbondingDelegation(
				delegation.DelegatorAddress,
				delegation.ValidatorAddress,
				sdk.NewCoin(bondDenom, entry.Balance),
				entry.CompletionTime,
				height,
				timestamp,
			))
		}
	}

	return unbondingDelegations, nil
}
