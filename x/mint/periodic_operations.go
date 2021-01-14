package mint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/utils"
)

// RegisterPeriodicOps returns the AdditionalOperation that periodically runs fetches from
// the LCD to make sure that constantly changing data are synced properly.
func RegisterPeriodicOps(scheduler *gocron.Scheduler, cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().Str("module", "mint").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateInflation(cp, db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateInflation fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func updateInflation(cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "mint").
		Str("operation", "inflation").
		Msg("getting inflation data")

	height, err := db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the inflation
	var inflation sdk.Dec
	endpoint := fmt.Sprintf("/mint/inflation?height=%d", height)
	_, err = cp.QueryLCDWithHeight(endpoint, &inflation)
	if err != nil {
		return err
	}

	return db.SaveInflation(inflation, height)
}
