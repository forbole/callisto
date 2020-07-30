package pricefeed

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/parse"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	operations "github.com/forbole/bdjuno/x/pricefeed/operations"
	"github.com/forbole/bdjuno/x/utils"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// PeriodicPriceFeedOperations returns the AdditionalOperation that periodically runs fetches from
// CoinGecko to make sure that constantly changing data are synced properly.
func PeriodicPriceFeedOperations(scheduler *gocron.Scheduler) parse.AdditionalOperation {
	log.Debug().Str("module", "PriceFeed").Msg("setting up periodic tasks")

	return func(_ config.Config, _ *codec.Codec, cp client.ClientProxy, db db.Database) error {
		bdDatabase, ok := db.(database.BigDipperDb)
		if !ok {
			log.Fatal().Str("module", "PriceFeed").Msg("given database instance is not a BigDipperDb")
		}

		// Fetch total supply of token in 30 seconds each
		if _, err := scheduler.Every(30).Second().StartImmediately().Do(func() {
			utils.WatchMethod(func() error { return operations.UpdatePrice(cp, bdDatabase) })
		}); err != nil {
			return err
		}
		return nil
	}
}
