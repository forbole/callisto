package pool

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
)

// UpdatePools refreshes protocol pools info in database
func (m *Module) UpdatePools(height int64) error {
	poolsList, err := m.source.Pools(height)
	if err != nil {
		return err
	}

	var pools []types.PoolList
	for _, pool := range poolsList {
		pools = append(pools, types.NewPoolList(
			pool.Id,
			pool.Data.Name,
			pool.Data.Runtime,
			pool.Data.Logo,
			pool.Data.Config,
			pool.Data.StartKey,
			pool.Data.CurrentKey,
			pool.Data.CurrentSummary,
			pool.Data.CurrentIndex,
			pool.Data.TotalBundles,
			pool.Data.UploadInterval,
			pool.Data.InflationShareWeight,
			pool.Data.MinDelegation,
			pool.Data.MaxBundleSize,
			pool.Data.Disabled,
			pool.Data.Protocol,
			pool.Data.UpgradePlan,
			pool.Data.CurrentStorageProviderId,
			pool.Data.CurrentCompressionId,
			height,
		))
	}

	err = m.db.SavePoolList(pools)
	if err != nil {
		return fmt.Errorf("error while saving pool list: %s", err)
	}

	return nil
}
