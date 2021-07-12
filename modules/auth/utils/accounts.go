package utils

import (
	"encoding/json"
	"context"


	"github.com/cosmos/cosmos-sdk/codec"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	codectypes"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	vesttypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/desmos-labs/juno/client"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"

	"github.com/forbole/bdjuno/types"
)

// GetGenesisAccounts parses the given appState and returns the genesis accounts
func GetGenesisAccounts(appState map[string]json.RawMessage, cdc codec.Marshaler) ([]types.Account, error) {
	log.Debug().Str("module", "auth").Msg("parsing genesis")

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
func GetAccounts(addresses []string, height int64,authClient authttypes.QueryClient) []codectypes.Any {
	log.Debug().Str("module", "auth").Str("operation", "accounts").Msg("getting accounts data")
	var accounts []codectypes.Any
	header := client.GetHeightRequestHeader(height)
	
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := authClient.Accounts(
			context.Background(),
			&authttypes.QueryAccountsRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 delegations at time
				},
			},
			header,
		)
		if err != nil {
			log.Error().Str("module", "auth").Err(err).Int64("height", height).
				Str("Auth","Get Account").Msg("error while getting accounts")
		}

		for _,account := range res.Accounts{
			accounts=append(accounts,*account)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		
	}
	return accounts
}


// UpdateAccounts takes the given addresses and for each one queries the chain
// retrieving the account data and stores it inside the database.
func UpdateAccounts(addresses []string, cdc codec.Marshaler, db *database.Db, height int64,authClient authttypes.QueryClient) error {
	accounts := GetAccounts(addresses,height,authClient)
	accountsList := make([]types.Account, len(accounts))

	for index,account :=range accounts{
		var accountI authttypes.BaseAccount
		err := cdc.UnpackAny(&account, accountI)
		if err != nil {
			return err
		}

		accountsList[index] = types.NewAccount(accountI.GetAddress().String())
	}
	
	return db.SaveAccounts(accountsList)
}
