package auth

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/desmos-labs/juno/client"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
)

// RefreshAccounts takes the given addresses and for each one queries the LCD
// retrieving the latest balance storing it inside the database.
func RefreshAccounts(addresses []string, height int64, cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().Str("module", "auth").Str("operation", "accounts").Msg("getting accounts data")

	// Get all the accounts information
	for _, address := range addresses {
		endpoint := fmt.Sprintf("/auth/accounts/%s?height=%d", address, height)

		var account exported.Account
		_, err := cp.QueryLCDWithHeight(endpoint, &account)
		if err != nil {
			log.Err(err).Str("module", "auth").Int64("height", height).Msg("error getting account")
			return err
		}

		err = db.SaveAccount(account)
		if err != nil {
			return err
		}

		err = db.SaveAccountBalance(account.GetAddress().String(), account.GetCoins(), height)
		if err != nil {
			return err
		}
	}

	return nil
}
