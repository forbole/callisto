package auth

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
)

// GetGenesisVestingAccounts parses the given appState and returns the genesis vesting accounts
func GetGenesisVestingAccounts(appState map[string]json.RawMessage, cdc codec.Codec) ([]exported.VestingAccount, error) {
	var authState authttypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[authttypes.ModuleName], &authState); err != nil {
		return nil, err
	}

	// Build vestingAccounts Array
	var vestingAccounts []exported.VestingAccount
	for _, account := range authState.Accounts {
		var accountI authttypes.AccountI
		err := cdc.UnpackAny(account, &accountI)
		if err != nil {
			return nil, err
		}

		vestingAccount, ok := accountI.(exported.VestingAccount)
		if !ok {
			continue
		}
		vestingAccounts = append(vestingAccounts, vestingAccount)
	}

	return vestingAccounts, nil
}
