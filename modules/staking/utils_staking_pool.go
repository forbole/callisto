package staking

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
	// sdk "github.com/cosmos/cosmos-sdk/types"
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

	var totalUnbondingTokens int64
	for _, validator := range validatorsList {
		unbondingDelegations, err := m.source.GetUnbondingDelegationsFromValidator(
			height,
			validator.GetOperator(),
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting unbonding delegations: %s", err)
		}

		for _, unbonding := range unbondingDelegations.UnbondingResponses {
			for _, entry := range unbonding.Entries {
				totalUnbondingTokens = totalUnbondingTokens + entry.Balance.Int64()
			}
		}
	}

	fmt.Printf("\n \n totalUnbondingTokens %v \n \n", totalUnbondingTokens)

	return types.NewPool(pool.BondedTokens, pool.NotBondedTokens, height), nil
}
