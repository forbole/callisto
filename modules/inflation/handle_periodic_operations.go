package inflation

import (
	"context"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/utils"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	inflationtypes "github.com/e-money/em-ledger/x/inflation/types"
)

// RegisterPeriodicOps returns the AdditionalOperation that periodically runs fetches from
// the LCD to make sure that constantly changing data are synced properly.
func RegisterPeriodicOps(scheduler *gocron.Scheduler, inflationClient inflationtypes.QueryClient, db *database.Db) error {
	log.Debug().Str("module", "inflation").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateInflation(inflationClient, db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateInflation fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func updateInflation(inflationClient inflationtypes.QueryClient, db *database.Db) error {
	log.Debug().
		Str("module", "inflation").
		Str("operation", "inflation").
		Msg("getting inflation data")

	// Get the inflation
	res, err := inflationClient.Inflation(
		context.Background(),
		&inflationtypes.QueryInflationRequest{},
	)
	if err != nil {
		return err
	}
	return db.SaveEmoneyInflation(res.State)
}
