package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/parse"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/utils"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// PeriodicAuthOperations returns the AdditionalOperation that periodically runs fetches from
// the LCD to make sure that constantly changing data are synced properly.
func PeriodicAuthOperations(scheduler *gocron.Scheduler) parse.AdditionalOperation {
	log.Debug().Str("module", "auth").Msg("setting up periodic tasks")

	return func(_ config.Config, _ *codec.Codec, cp client.ClientProxy, db db.Database) error {
		bdDatabase, ok := db.(database.BigDipperDb)
		if !ok {
			log.Fatal().Msg("given database instance is not a BigDipperDb")
		}

		// Setup a cron job to run every midnight
		if _, err := scheduler.Every(1).Day().At("00:00").StartImmediately().Do(func() {
			utils.WatchMethod(func() error { return updateAccounts(cp, bdDatabase) })
		}); err != nil {
			return err
		}

		return nil
	}
}
