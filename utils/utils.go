package utils

import "github.com/forbole/bdjuno/v2/types"

// RemoveDuplicateValues removes the duplicated values from the given slice
func RemoveDuplicateValues(slice []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// RemoveDuplicateAccountBalance removes the duplicated values of account balance from the given slice
func RemoveDuplicateAccountBalance(intSlice []types.AccountBalance) []types.AccountBalance {
	var list []types.AccountBalance
	for _, entry := range intSlice {
		if len(list) == 0 {
			list = append(list, entry)
		}
		if len(list) > 0 {
			for _, listElement := range list {
				if entry.Address != listElement.Address {
					list = append(list, entry)
				}
			}
		}
	}
	return list
}
