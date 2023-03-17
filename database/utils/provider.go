package utils

import (
	"github.com/forbole/bdjuno/v4/types"
)

func SplitProviders(providers []*types.Provider, paramsNumber int) [][]*types.Provider {
	maxBalancesPerSlice := maxPostgreSQLParams / paramsNumber
	slices := make([][]*types.Provider, len(providers)/maxBalancesPerSlice+1)

	sliceIndex := 0
	for index, provider := range providers {
		slices[sliceIndex] = append(slices[sliceIndex], provider)

		if index > 0 && index%(maxBalancesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}
