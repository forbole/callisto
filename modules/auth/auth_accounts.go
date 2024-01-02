package auth

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/types"
)

// GetGenesisAccounts parses the given appState and returns the genesis accounts
func GetGenesisAccounts(appState map[string]json.RawMessage, cdc codec.Codec) ([]types.Account, error) {
	var authState authttypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[authttypes.ModuleName], &authState); err != nil {
		return nil, err
	}

	// Store the accounts
	accounts := make([]types.Account, len(authState.Accounts))
	for index, account := range authState.Accounts {
		var accountI authttypes.AccountI
		err := cdc.UnpackAny(account, &accountI)
		if err != nil {
			return nil, err
		}

		accounts[index] = types.NewAccount(accountI.GetAddress().String())
	}

	return accounts, nil
}

// --------------------------------------------------------------------------------------------------------------------

// GetAccounts returns the account data for the given addresses
func GetAccounts(height int64, addresses []string) []types.Account {
	log.Debug().Str("module", "auth").Str("operation", "accounts").Msg("getting accounts data")

	// Get all the accounts information
	var accounts = make([]types.Account, len(addresses))
	for index, address := range addresses {
		accounts[index] = types.NewAccount(address)
	}

	return accounts
}

// RefreshAccounts takes the given addresses and for each one queries the chain
// retrieving the account data and stores it inside the database.
func (m *Module) RefreshAccounts(height int64, addresses []string) error {
	if len(addresses) == 0 {
		return nil
	}
	accounts := GetAccounts(height, addresses)
	return m.db.SaveAccounts(accounts)
}
