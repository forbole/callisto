package distribution

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/utils"
)

// PeriodicStakingOperations returns the AdditionalOperation that periodically runs fetches from
// the LCD to make sure that constantly changing data are synced properly.
func RegisterPeriodicOps(
	scheduler *gocron.Scheduler, distrClient distrtypes.QueryClient, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "distribution").Msg("setting up periodic tasks")

	// Fetch community pool in 30 seconds each
	if _, err := scheduler.Every(30).Second().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateCommunityPool(distrClient, db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateCommunityPool fetch total amount of coins in the system from RPC and store it into database
func updateCommunityPool(distrClient distrtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "distribution").Str("operation", "community pool").
		Msg("getting community pool")

	height, err := db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	res, err := distrClient.CommunityPool(context.Background(), &distrtypes.QueryCommunityPoolRequest{})
	if err != nil {
		return err
	}

	log.Debug().Str("module", "distribution").Str("operation", "community pool").
		Msg("saving community pool")

	// Store the signing infos into the database
	return db.SaveCommunityPool(res.Pool, height)
}
