package staking

import (
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/common"
	"github.com/forbole/bdjuno/x/utils"
)

// RegisterPeriodicOps registers the additional operations that periodically run
func RegisterPeriodicOps(
	scheduler *gocron.Scheduler, stakingClient staking.QueryClient, cdc codec.Marshaler, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "staking").Msg("setting up periodic tasks")

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
	stakingClient staking.QueryClient, cdc codec.Marshaler, db *database.BigDipperDb,
) error {
	height, err := db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the params
	params, err := common.UpdateParams(height, stakingClient, db)
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
		common.UpdateValidatorsDelegations(height, validators, stakingClient, db)
		wg.Done()
	}()

	// Update the unbonding delegations
	wg.Add(1)
	go func() {
		common.UpdateValidatorsUnbondingDelegations(height, params.BondDenom, validators, stakingClient, db)
		wg.Done()
	}()

	// Update the redelegations
	wg.Add(1)
	go func() {
		common.UpdateValidatorsRedelegations(height, params.BondDenom, validators, stakingClient, db)
		wg.Done()
	}()

	wg.Wait()
	return nil
}
