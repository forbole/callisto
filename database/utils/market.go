package utils

import (
	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"
)

func SplitLeases(leases []markettypes.QueryLeaseResponse, paramsNumber int) [][]markettypes.QueryLeaseResponse {
	maxLeasesPerSlice := maxPostgreSQLParams / paramsNumber
	slices := make([][]markettypes.QueryLeaseResponse, len(leases)/maxLeasesPerSlice+1)

	sliceIndex := 0
	for index, lease := range leases {
		slices[sliceIndex] = append(slices[sliceIndex], lease)

		if index > 0 && index%(maxLeasesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}
