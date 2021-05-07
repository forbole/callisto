package distribution

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/modules/common/utils"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	bdistrcommon "github.com/forbole/bdjuno/modules/bigdipper/distribution/common"
)

// RegisterPeriodicOps registers the additional operations that periodically run
func RegisterPeriodicOps(
	scheduler *gocron.Scheduler, distrClient distrtypes.QueryClient, db *bigdipperdb.Db,
) error {
	log.Debug().Str("module", "distribution").Msg("setting up periodic tasks")

	// Update the community pool every 1 hour
	if _, err := scheduler.Every(1).Hour().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return getLatestCommunityPool(distrClient, db) })
	}); err != nil {
		return err
	}

	return nil
}

// getLatestCommunityPool gets the latest community pool from the chain and stores inside the database
func getLatestCommunityPool(distrClient distrtypes.QueryClient, db *bigdipperdb.Db) error {
	height, err := db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	return bdistrcommon.UpdateCommunityPool(height, distrClient, db)
}
