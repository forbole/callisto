package utils

import "github.com/forbole/bdjuno/types"

func SplitBalances(balances []types.AccountBalance, paramsNumber int) [][]types.AccountBalance {
	maxPostgreSQLParams := 65535
	maxBalancesPerSlice := maxPostgreSQLParams / paramsNumber

	slices := make([][]types.AccountBalance, len(balances)/maxBalancesPerSlice+1)

	sliceIndex := 0
	for index, balance := range balances {
		slices[sliceIndex] = append(slices[sliceIndex], balance)

		if index > 0 && index%(maxBalancesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}
