package auth

import (
	"context"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// FastSync downloads the x/auth state at the given height, and stores it inside the database
func FastSync(height int64, client authttypes.QueryClient, db *database.Db) error {
	err := updateAccounts(height, client, db)
	if err != nil {
		return err
	}

	return nil
}

// updateAccounts downloads all the accounts at the given height, and stores them inside the database
func updateAccounts(height int64, authClient authttypes.QueryClient, db *database.Db) error {
	header := client.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := authClient.Accounts(
			context.Background(),
			&authttypes.QueryAccountsRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 accounts at time
				},
			},
			header,
		)
		if err != nil {
			return err
		}

		var accounts = make([]types.Account, len(res.Accounts))
		for index, acc := range res.Accounts {
			accounts[index] = types.NewAccount(
				acc.GetCachedValue().(authttypes.AccountI).GetAddress().String(),
			)
		}

		err = db.SaveAccounts(accounts)
		if err != nil {
			return err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
	}

	return nil
}
