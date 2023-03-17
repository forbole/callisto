package utils

import "github.com/forbole/bdjuno/v4/types"

func SplitWasmContracts(contracts []types.WasmContract, paramsNumber int) [][]types.WasmContract {
	maxBalancesPerSlice := maxPostgreSQLParams / paramsNumber
	slices := make([][]types.WasmContract, len(contracts)/maxBalancesPerSlice+1)

	sliceIndex := 0
	for index, contract := range contracts {
		slices[sliceIndex] = append(slices[sliceIndex], contract)

		if index > 0 && index%(maxBalancesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}

func SplitWasmExecuteContracts(executeContracts []types.WasmExecuteContract, paramsNumber int) [][]types.WasmExecuteContract {
	maxBalancesPerSlice := maxPostgreSQLParams / paramsNumber
	slices := make([][]types.WasmExecuteContract, len(executeContracts)/maxBalancesPerSlice+1)

	sliceIndex := 0
	for index, executeContract := range executeContracts {
		slices[sliceIndex] = append(slices[sliceIndex], executeContract)

		if index > 0 && index%(maxBalancesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}
