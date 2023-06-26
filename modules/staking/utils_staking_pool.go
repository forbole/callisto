package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/forbole/bdjuno/v4/types"
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

	var unbondingTokens = sdk.NewInt(0)

	for _, validator := range validatorsList {
		// get list of all unbonding delegations for each validator
		unbondingDelegations := m.getTotalUnbondingDelegationsFromValidator(height, validator.GetOperator())
		if len(unbondingDelegations) > 0 {
			// calculate total value of unbonding tokens
			for _, unbonding := range unbondingDelegations {
				for _, entry := range unbonding.Entries {
					unbondingTokens = unbondingTokens.Add(entry.Balance)
				}
			}
		}

	}

	// calculate total value of staked tokens that are not bonded
	stakedNotBondedTokens := pool.NotBondedTokens.Sub(unbondingTokens)

	return types.NewPool(pool.BondedTokens, pool.NotBondedTokens, unbondingTokens, stakedNotBondedTokens, height), nil
}

func (m *Module) GetStakingPoolSnapshot(height int64) (*types.PoolSnapshot, error) {
	pool, err := m.source.GetPool(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting staking pool snapshot: %s", err)
	}

	return types.NewPoolSnapshot(pool.BondedTokens, pool.NotBondedTokens, height), nil
}

func (m *Module) getTotalUnbondingDelegationsFromValidator(height int64, valOperatorAddress string) []stakingtypes.UnbondingDelegation {
	var unbondingDelegations []stakingtypes.UnbondingDelegation
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.source.GetUnbondingDelegationsFromValidator(height,
			valOperatorAddress,
			&query.PageRequest{Key: nextKey},
		)
		if err != nil {
			return []stakingtypes.UnbondingDelegation{}
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		unbondingDelegations = append(unbondingDelegations, res.UnbondingResponses...)
	}
	return unbondingDelegations
}
