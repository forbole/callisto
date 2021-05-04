package auth

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/utils"
)

// FastSync downloads the x/auth state at the given height, and stores it inside the database
func FastSync(height int64, client authtypes.QueryClient, db *database.BigDipperDb) error {
	err := updateAccounts(height, client, db)
	if err != nil {
		return err
	}

	return nil
}

// updateAccounts downloads all the accounts at the given height, and stores them inside the database
func updateAccounts(height int64, client authtypes.QueryClient, db *database.BigDipperDb) error {
	header := utils.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := client.Accounts(
			context.Background(),
			&authtypes.QueryAccountsRequest{
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

		var accounts = make([]authtypes.AccountI, len(res.Accounts))
		for index, acc := range res.Accounts {
			accounts[index] = acc.GetCachedValue().(authtypes.AccountI)
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
