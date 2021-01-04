package distribution

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/utils"
)

// PeriodicStakingOperations returns the AdditionalOperation that periodically runs fetches from
// the LCD to make sure that constantly changing data are synced properly.
func RegisterPeriodicOps(scheduler *gocron.Scheduler, cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().Str("module", "distribution").Msg("setting up periodic tasks")

	// Fetch community pool in 30 seconds each
	if _, err := scheduler.Every(30).Second().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateCommunityPool(cp, db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateCommunityPool fetch total amount of coins in the system from RPC and store it into database
func updateCommunityPool(cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "distribution").
		Str("operation", "community pool").
		Msg("getting community pool")

	// Get the community pool
	var s sdk.DecCoins
	height, err := cp.QueryLCDWithHeight("/distribution/community_pool", &s)
	if err != nil {
		return err
	}

	// Store the signing infos into the database
	return db.SaveCommunityPool(s, height)
}
