package utils

import (
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"
)

func SplitProviders(providers []providertypes.Provider, paramsNumber int) [][]providertypes.Provider {
	maxProvidersPerSlice := maxPostgreSQLParams / paramsNumber
	slices := make([][]providertypes.Provider, len(providers)/maxProvidersPerSlice+1)

	sliceIndex := 0
	for index, provider := range providers {
		slices[sliceIndex] = append(slices[sliceIndex], provider)

		if index > 0 && index%(maxProvidersPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}
