package authority

import (
	"context"

	authoritytypes "github.com/e-money/em-ledger/x/authority/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/utils"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOps returns the updateMinGasPrices that periodically fetches latest min gas prices
func RegisterPeriodicOps(
	scheduler *gocron.Scheduler,
	authorityclient authoritytypes.QueryClient,
	db *database.Db,
) error {
	log.Debug().Str("module", "authority").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateMinGasPrices(authorityclient, db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateMinGasPrices fetches from the REST APIs the latest value for the
// minimum gas prices, and saves it inside the database.
func updateMinGasPrices(authorityclient authoritytypes.QueryClient, db *database.Db) error {
	log.Debug().
		Str("module", "authority").
		Str("operation", "gas pirces").
		Msg("querying minimum gas pirces")

	height, err := db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the min gas prices
	res, err := authorityclient.GasPrices(
		context.Background(),
		&authoritytypes.QueryGasPricesRequest{},
	)
	if err != nil {
		return err
	}

	return db.UpdateEMoneyGasPrices(res.GetMinGasPrices(), height)
}
