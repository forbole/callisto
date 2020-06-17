package auth

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// RefreshAccounts takes the given addresses and for each one queries the LCD
// retrieving the latest balance storing it inside the database.
func RefreshAccounts(
	addresses []sdk.AccAddress, height int64, timestamp time.Time, cp client.ClientProxy, db database.BigDipperDb,
) error {
	// Get all the accounts information
	var accounts []exported.Account
	for _, address := range addresses {
		endpoint := fmt.Sprintf("/auth/accounts/%s?height=%d", address.String(), height)

		var account exported.Account
		_, err := cp.QueryLCDWithHeight(endpoint, &account)
		if err != nil {
			return err
		}

		accounts = append(accounts, account)
	}

	return db.SaveAccounts(accounts, height, timestamp)
}

// updateAccounts gets all the accounts stored inside the database, and refreshes their
// balances by fetching the LCD endpoint.
func updateAccounts(cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().Str("module_name", "auth").Msg("updating accounts")

	var block tmctypes.ResultBlock
	err := cp.QueryLCD("/blocks/latest", &block)
	if err != nil {
		return err
	}

	addresses, err := db.GetAccounts()
	if err != nil {
		return err
	}

	return RefreshAccounts(addresses, block.Block.Height, block.Block.Time, cp, db)
}
