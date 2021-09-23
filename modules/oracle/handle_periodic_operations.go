package oracle

import (
	"context"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/utils"
	"github.com/forbole/bdjuno/types"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	oracletypes "github.com/bandprotocol/chain/v2/x/oracle/types"
)

// RegisterPeriodicOps returns the AdditionalOperation that periodically runs fetches from
// the LCD to make sure that constantly changing data are synced properly.
func RegisterPeriodicOps(scheduler *gocron.Scheduler, oracleClient oracletypes.QueryClient, db *database.Db) error {
	log.Debug().Str("module", "oracle").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateOracleParams(oracleClient, db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateOracleParams fetches from the REST APIs the latest value for
// oracle params, and saves it inside the database.
func updateOracleParams(oracleClient oracletypes.QueryClient, db *database.Db) error {
	log.Debug().
		Str("module", "oracle").
		Msg("getting inflation data")

	height, err := db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the inflation
	res, err := oracleClient.Params(
		context.Background(),
		&oracletypes.QueryParamsRequest{},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	paramsList := types.NewOracleParams(res.Params, height)

	return db.SaveOracleParams(paramsList, height)
}
