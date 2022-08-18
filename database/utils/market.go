package utils

import (
	"github.com/forbole/bdjuno/v3/types"
)

func SplitLeases(leases []*types.MarketLease, paramsNumber int) [][]*types.MarketLease {
	maxBalancesPerSlice := maxPostgreSQLParams / paramsNumber
	slices := make([][]*types.MarketLease, len(leases)/maxBalancesPerSlice+1)

	sliceIndex := 0
	for i, lease := range leases {
		slices[sliceIndex] = append(slices[sliceIndex], lease)

		if i > 0 && i%(maxBalancesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}
