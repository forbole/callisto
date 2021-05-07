package mint

import (
	"context"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	utils2 "github.com/forbole/bdjuno/modules/common/utils"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// RegisterPeriodicOps returns the AdditionalOperation that periodically runs fetches from
// the LCD to make sure that constantly changing data are synced properly.
func RegisterPeriodicOps(scheduler *gocron.Scheduler, minClient minttypes.QueryClient, db *bigdipperdb.Db) error {
	log.Debug().Str("module", "mint").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").StartImmediately().Do(func() {
		utils2.WatchMethod(func() error { return updateInflation(minClient, db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateInflation fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func updateInflation(mintClient minttypes.QueryClient, db *bigdipperdb.Db) error {
	log.Debug().
		Str("module", "mint").
		Str("operation", "inflation").
		Msg("getting inflation data")

	height, err := db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the inflation
	res, err := mintClient.Inflation(
		context.Background(),
		&minttypes.QueryInflationRequest{},
		utils2.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	return db.SaveInflation(res.Inflation, height)
}
