package utils

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/types"
)

// GetGenesisAccounts parses the given appState and returns the genesis accounts
func GetGenesisAccounts(appState map[string]json.RawMessage, cdc codec.Marshaler) ([]types.Account, error) {
	log.Debug().Str("module", "auth").Msg("parsing genesis")

	var authState authtypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[authtypes.ModuleName], &authState); err != nil {
		return nil, err
	}

	// Store the accounts
	accounts := make([]types.Account, len(authState.Accounts))
	for index, account := range authState.Accounts {
		var accountI authtypes.AccountI
		err := cdc.UnpackAny(account, &accountI)
		if err != nil {
			return nil, err
		}

		accounts[index] = types.NewAccount(accountI.GetAddress().String())
	}

	return accounts, nil
}

// GetAccounts returns the account data for the given addresses
func GetAccounts(addresses []string) []types.Account {
	log.Debug().Str("module", "auth").Str("operation", "accounts").Msg("getting accounts data")

	// Get all the accounts information
	var accounts = make([]types.Account, len(addresses))
	for index, address := range addresses {
		accounts[index] = types.NewAccount(address)
	}

	return accounts
}
