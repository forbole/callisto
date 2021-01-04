package supply

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/utils"
)

func RegisterPeriodicOps(scheduler *gocron.Scheduler, cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().Str("module", "supply").Msg("setting up periodic tasks")

	// Fetch total supply of token in 30 seconds each
	_, err := scheduler.Every(30).Second().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateTotalTokenSupply(cp, db) })
	})
	if err != nil {
		return err
	}

	return nil
}

// updateTotalTokenSupply fetch total amount of coins in the system from RPC and store it into database
func updateTotalTokenSupply(cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", "tokens").
		Msg("getting total token supply")

	// Get the total supply
	var s sdk.Coins
	height, err := cp.QueryLCDWithHeight("/supply/total", &s)
	if err != nil {
		return err
	}

	// Store the signing infos into the database
	return db.SaveSupplyToken(s, height)
}
