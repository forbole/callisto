package utils

import (
	"encoding/json"
	"context"


	"github.com/cosmos/cosmos-sdk/codec"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/types/query"

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

		accounts[index] = types.NewAccount(accountI.GetAddress().String(),accountI)
	}

	return accounts, nil
}

// --------------------------------------------------------------------------------------------------------------------

// GetAllAccounts returns all account data
func GetAllAccounts(addresses []string,cdc codec.Marshaler,height int64,authClient authttypes.QueryClient) ([]types.Account,error) {
	log.Debug().Str("module", "auth").Str("operation", "accounts").Msg("getting accounts data")
	var accounts []types.Account
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
			var accountI authttypes.AccountI
			err := cdc.UnpackAny(account, &accountI)
			if err!=nil{
				return nil,err
			}
			accounts=append(accounts, types.NewAccount(accountI.GetAddress().String(),accountI))

		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		
	}


	return accounts,nil
}

// GetAccounts returns the account data for the given addresses
func GetAccounts(addresses []string,cdc codec.Marshaler,height int64,authClient authttypes.QueryClient) ([]types.Account,error) {
	log.Debug().Str("module", "auth").Str("operation", "accounts").Msg("getting accounts data")
	var accounts []types.Account
	header := client.GetHeightRequestHeader(height)
	
	for _,address:=range addresses{
		res, err := authClient.Account(
			context.Background(),
			&authttypes.QueryAccountRequest{
				Address: address,
			},
			header,
		)
		if err != nil {
			log.Error().Str("module", "auth").Err(err).Int64("height", height).
				Str("Auth","Get Account").Msg("error while getting accounts")
		}
			var accountI authttypes.AccountI
			err = cdc.UnpackAny(res.Account, &accountI)
			if err!=nil{
				return nil,err
			}
			accounts=append(accounts, types.NewAccount(accountI.GetAddress().String(),accountI))

}


	return accounts,nil
}


// UpdateAccounts takes the given addresses and for each one queries the chain
// retrieving the account data and stores it inside the database.
func UpdateAccounts(addresses []string, cdc codec.Marshaler, db *database.Db, height int64,authClient authttypes.QueryClient) error {
	accounts,err := GetAccounts(addresses,cdc,height,authClient)
	if err!=nil{
		return err
	}

	return db.SaveAccounts(accounts)
}
