package utils

import "github.com/forbole/bdjuno/v2/types"

func SplitDelegations(delegations []types.Delegation, paramsNumber int) [][]types.Delegation {
	maxBalancesPerSlice := maxPostgreSQLParams / paramsNumber
	slices := make([][]types.Delegation, len(delegations)/maxBalancesPerSlice+1)

	sliceIndex := 0
	for index, delegation := range delegations {
		slices[sliceIndex] = append(slices[sliceIndex], delegation)

		if index > 0 && index%(maxBalancesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}

func SplitRedelegations(redelegations []types.Redelegation, paramsNumber int) [][]types.Redelegation {
	maxBalancesPerSlice := maxPostgreSQLParams / paramsNumber
	slices := make([][]types.Redelegation, len(redelegations)/maxBalancesPerSlice+1)

	sliceIndex := 0
	for index, redelegation := range redelegations {
		slices[sliceIndex] = append(slices[sliceIndex], redelegation)

		if index > 0 && index%(maxBalancesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}

func SplitUnbondingDelegations(unbondingDelegations []types.UnbondingDelegation, paramsNumber int) [][]types.UnbondingDelegation {
	maxBalancesPerSlice := maxPostgreSQLParams / paramsNumber
	slices := make([][]types.UnbondingDelegation, len(unbondingDelegations)/maxBalancesPerSlice+1)

	sliceIndex := 0
	for index, unbondingDelegation := range unbondingDelegations {
		slices[sliceIndex] = append(slices[sliceIndex], unbondingDelegation)

		if index > 0 && index%(maxBalancesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}
