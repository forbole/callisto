package consensus

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/parse"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/consensus/operations"
	"github.com/forbole/bdjuno/x/utils"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// PeriodicConcensusOperations returns the AdditionalOperation that periodically runs fetches from
// the LCD to make sure that constantly changing data are synced properly.
func PeriodicConcensusOperations(scheduler *gocron.Scheduler) parse.AdditionalOperation {
	log.Debug().Str("module", "consensus").Msg("setting up periodic tasks")

	return func(_ config.Config, _ *codec.Codec, cp client.ClientProxy, db db.Database) error {
		bdDatabase, ok := db.(database.BigDipperDb)
		if !ok {
			log.Fatal().Str("module", "consensus").Msg("given database instance is not a BigDipperDb")
		}

		// Setup a cron job to run every 15 seconds
		if _, err := scheduler.Every(1).Minute().At("00:00").Do(func() {
			utils.WatchMethod(func() error { return operations.UpdateBlockTimeInMinute(cp, bdDatabase) })
		}); err != nil {
			return err
		}
		/*
			if _, err := scheduler.Every(1).Minute().At("00:00").StartImmediately().Do(func() {
				utils.WatchMethod(func() error { return operations.UpdateValidatorVotingPower(cp, bdDatabase) })
			}); err != nil {
				return err
			}
			if _, err := scheduler.Every(1).Hour().At("00:00").StartImmediately().Do(func() {
				utils.WatchMethod(func() error { return operations.UpdateValidatorVotingPower(cp, bdDatabase) })
			}); err != nil {
				return err
			}
		*/
		return nil
	}
}
