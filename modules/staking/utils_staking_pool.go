package staking

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
)

func (m *Module) GetStakingPool(height int64) (*types.Pool, error) {
	pool, err := m.source.GetPool(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting staking pool: %s", err)
	}

	return types.NewPool(pool.BondedTokens, pool.NotBondedTokens, height), nil
}
