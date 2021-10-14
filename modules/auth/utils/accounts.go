package utils

import (
	"encoding/json"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"

	"github.com/forbole/bdjuno/types"
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

// GetGenesisVestingAccounts parses the given appState and returns the desired genesis vesting account details
func GetGenesisVestingAccounts(appState map[string]json.RawMessage, cdc codec.Codec) ([]types.VestingAccount, error) {
	log.Debug().Str("module", "auth (vesting)").Msg("parsing genesis")

	var authState Accounts
	if err := json.Unmarshal(appState[authttypes.ModuleName], &authState); err != nil {
		return nil, err
	}

	// Store the accounts
	accounts := []types.VestingAccount{}
	for _, account := range authState.Accounts {
		accountType := strings.Split(account.AccountType, ".")[3]

		if accountType == "PeriodicVestingAccount" {
			// Prepare vesting periods array for creating new Account instance
			vestingPeriods := make([]types.VestingPeriod, len(account.VestingPeriods))
			for index, period := range account.VestingPeriods {
				vestingPeriods[index] = types.NewVestingPeriod(period.Length, period.Amount)
			}

			// Create new Account instance
			Account := types.NewVestingAccount(
				account.BaseVestingAccount.BaseAccount.Address,
				account.BaseVestingAccount.OriginalVesting,
				account.BaseVestingAccount.EndTime,
				account.StartTime,
				vestingPeriods,
			)
			accounts = append(accounts, Account)
		}
	}

	return accounts, nil
}

// --------------------------------------------------------------------------------------------------------------------

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

// UpdateAccounts takes the given addresses and for each one queries the chain
// retrieving the account data and stores it inside the database.
func UpdateAccounts(addresses []string, db *database.Db) error {
	accounts := GetAccounts(addresses)
	return db.SaveAccounts(accounts)
}
