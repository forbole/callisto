package staking

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

// PeriodicStakingOperations returns the AdditionalOperation that periodically runs fetches from
// the LCD to make sure that constantly changing data are synced properly.
func PeriodicStakingOperations(scheduler *gocron.Scheduler) parse.AdditionalOperation {
	log.Debug().Str("module", "staking").Msg("setting up periodic tasks")

	return func(_ config.Config, _ *codec.Codec, cp client.ClientProxy, db db.Database) error {
		bdDatabase, ok := db.(database.BigDipperDb)
		if !ok {
			log.Fatal().Str("module", "staking").Msg("given database instance is not a BigDipperDb")
		}

		// Setup a cron job to run every 15 seconds
		if _, err := scheduler.Every(15).Second().StartImmediately().Do(func() {
			utils.WatchMethod(func() error { return updateStakingPool(cp, bdDatabase) })
			utils.WatchMethod(func() error { return updateValidatorsUptime(cp, bdDatabase) })
		}); err != nil {
			return err
		}

		// Setup a cron job to run every midnight
		if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
			utils.WatchMethod(func() error { return updateDelegations(cp, bdDatabase) })
			utils.WatchMethod(func() error { return updateStakingPool(cp, bdDatabase) })
		}); err != nil {
			return err
		}

		return nil
	}
}
