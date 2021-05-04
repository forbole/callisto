package common

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/x/utils"

	"github.com/forbole/bdjuno/database"
)

// UpdateAccounts takes the given addresses and for each one queries the chain
// retrieving the account data and stores it inside the database.
func UpdateAccounts(
	addresses []string, height int64, authClient authtypes.QueryClient,
	marshaler codec.Marshaler, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "auth").Str("operation", "accounts").Msg("getting accounts data")
	header := utils.GetHeightRequestHeader(height)

	// Get all the accounts information
	var accounts []authtypes.AccountI
	for _, address := range addresses {
		accRes, err := authClient.Account(
			context.Background(),
			&authtypes.QueryAccountRequest{Address: address},
			header,
		)
		if err != nil {
			return err
		}

		var account authtypes.AccountI
		err = marshaler.UnpackAny(accRes.Account, &account)
		if err != nil {
			return err
		}

		accounts = append(accounts, account)
	}

	return db.SaveAccounts(accounts)
}
