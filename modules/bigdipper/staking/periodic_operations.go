package staking

import (
	"sync"

	staking2 "github.com/forbole/bdjuno/modules/common/staking"
	utils2 "github.com/forbole/bdjuno/modules/common/utils"

	"github.com/cosmos/cosmos-sdk/codec"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	bstakingcommon "github.com/forbole/bdjuno/modules/bigdipper/staking/common"
)

// RegisterPeriodicOps registers the additional operations that periodically run
func RegisterPeriodicOps(
	scheduler *gocron.Scheduler, stakingClient staking.QueryClient, cdc codec.Marshaler, db *bigdipperdb.Db,
) error {
	log.Debug().Str("module", "staking").Msg("setting up periodic tasks")

	// Update the validator delegations every 1 hour
	if _, err := scheduler.Every(1).Hour().StartImmediately().Do(func() {
		utils2.WatchMethod(func() error { return updateValidatorsDelegations(stakingClient, cdc, db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateValidatorsDelegations updates the current validators set and all their delegations, unbonding delegations
// and redelegations
func updateValidatorsDelegations(
	stakingClient staking.QueryClient, cdc codec.Marshaler, db *bigdipperdb.Db,
) error {
	height, err := db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the params
	params, err := bstakingcommon.UpdateParams(height, stakingClient, db)
	if err != nil {
		return err
	}

	validators, err := updateValidators(height, stakingClient, cdc, db)
	if err != nil {
		return err
	}

	// Update the delegations
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		staking2.UpdateValidatorsDelegations(height, validators, stakingClient, db)
		wg.Done()
	}()

	// Update the unbonding delegations
	wg.Add(1)
	go func() {
		bstakingcommon.UpdateValidatorsUnbondingDelegations(height, params.BondDenom, validators, stakingClient, db)
		wg.Done()
	}()

	// Update the redelegations
	wg.Add(1)
	go func() {
		bstakingcommon.UpdateValidatorsRedelegations(height, params.BondDenom, validators, stakingClient, db)
		wg.Done()
	}()

	wg.Wait()
	return nil
}
