package auth

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/desmos-labs/juno/client"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
)

// RefreshAccounts takes the given addresses and for each one queries the LCD
// retrieving the latest balance storing it inside the database.
func RefreshAccounts(
	addresses []sdk.AccAddress, height int64, timestamp time.Time, cp *client.Proxy, db *database.BigDipperDb,
) error {
	log.Debug().
		Str("module", "auth").
		Str("operation", "accounts").
		Msg("getting accounts data")

	// Get all the accounts information
	var accounts []exported.Account
	for _, address := range addresses {
		endpoint := fmt.Sprintf("/auth/accounts/%s?height=%d", address.String(), height)

		var account exported.Account
		_, err := cp.QueryLCDWithHeight(endpoint, &account)
		if err != nil {
			log.Err(err).Str("module", "auth").Int64("height", height).Msg("error getting account")
			return err
		}

		accounts = append(accounts, account)
	}

	log.Debug().
		Str("module", "auth").
		Str("operation", "accounts").
		Msg("saving accounts data")
	return db.SaveAccounts(accounts, height, timestamp)
}
