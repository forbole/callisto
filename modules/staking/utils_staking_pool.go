package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v3/types"
)

func (m *Module) GetStakingPool(height int64) (*types.Pool, error) {
	pool, err := m.source.GetPool(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting staking pool: %s", err)
	}

	validatorsList, err := m.db.GetValidators()
	if err != nil {
		return nil, fmt.Errorf("error while getting validators list: %s", err)
	}

	var unbondingTokens int64
	for _, validator := range validatorsList {
		unbondingDelegations, err := m.source.GetUnbondingDelegationsFromValidator(
			height,
			validator.GetOperator(),
			nil,
		)

		if err != nil {
			// Do nothing
			continue
		}

		// calculate total value of unbonding tokens
		for _, unbonding := range unbondingDelegations.UnbondingResponses {
			for _, entry := range unbonding.Entries {
				// add to total unbonding value
				unbondingTokens += entry.Balance.Int64()
			}
		}
	}

	// calculate total value of staked tokens that are not bonded
	stakedNotBondedTokens := pool.NotBondedTokens.Int64() - unbondingTokens

	return types.NewPool(pool.BondedTokens, pool.NotBondedTokens, sdk.NewInt(unbondingTokens), sdk.NewInt(stakedNotBondedTokens), height), nil
}
