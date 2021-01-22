package auth

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/x/utils"

	"github.com/forbole/bdjuno/database"
)

// RefreshAccounts takes the given addresses and for each one queries the chain
// retrieving the latest balance and stores it inside the database.
func RefreshAccounts(
	addresses []string, height int64,
	authClient authtypes.QueryClient, bankClient banktypes.QueryClient, marshaler codec.Marshaler,
	db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "auth").Str("operation", "accounts").Msg("getting accounts data")
	header := utils.GetHeightRequestHeader(height)

	// Get all the accounts information
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

		balRes, err := bankClient.AllBalances(
			context.Background(),
			&banktypes.QueryAllBalancesRequest{Address: address},
			header,
		)
		if err != nil {
			return err
		}

		err = db.SaveAccount(account)
		if err != nil {
			return err
		}

		err = db.SaveAccountBalance(account.GetAddress().String(), balRes.Balances, height)
		if err != nil {
			return err
		}
	}

	return nil
}
