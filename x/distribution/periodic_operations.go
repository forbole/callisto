package distribution

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/x/distribution/common"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/utils"
)

// RegisterPeriodicOps registers the additional operations that periodically run
func RegisterPeriodicOps(
	scheduler *gocron.Scheduler, distrClient distrtypes.QueryClient, db *database.BigDipperDb,
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
func getLatestCommunityPool(distrClient distrtypes.QueryClient, db *database.BigDipperDb) error {
	height, err := db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	return common.UpdateCommunityPool(height, distrClient, db)
}
