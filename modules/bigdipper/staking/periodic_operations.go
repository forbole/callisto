package staking

import (
	"sync"

	"github.com/forbole/bdjuno/modules/common/staking"
	"github.com/forbole/bdjuno/modules/common/utils"

	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	bstakingcommon "github.com/forbole/bdjuno/modules/bigdipper/staking/common"
)

// RegisterPeriodicOps registers the additional operations that periodically run
func RegisterPeriodicOps(
	scheduler *gocron.Scheduler, stakingClient stakingtypes.QueryClient, cdc codec.Marshaler, db *bigdipperdb.Db,
) error {
	log.Debug().Str("module", "stakingtypes").Msg("setting up periodic tasks")

	// Update the validator delegations every 1 hour
	if _, err := scheduler.Every(1).Hour().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateValidatorsDelegations(stakingClient, cdc, db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateValidatorsDelegations updates the current validators set and all their delegations, unbonding delegations
// and redelegations
func updateValidatorsDelegations(
	stakingClient stakingtypes.QueryClient, cdc codec.Marshaler, db *bigdipperdb.Db,
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

	validators, err := staking.UpdateValidators(height, stakingClient, cdc, db)
	if err != nil {
		return err
	}

	// Update the delegations
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		staking.UpdateValidatorsDelegations(height, validators, stakingClient, db)
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
