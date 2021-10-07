package staking

import (
	"fmt"
	"sync"

	"github.com/forbole/bdjuno/database"
	stakingutils "github.com/forbole/bdjuno/modules/staking/utils"
	"github.com/forbole/bdjuno/modules/utils"

	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOps registers the additional utils that periodically run
func RegisterPeriodicOps(
	scheduler *gocron.Scheduler, stakingClient stakingtypes.QueryClient, cdc codec.Marshaler, db *database.Db,
) error {
	log.Debug().Str("module", "stakingtypes").Msg("setting up periodic tasks")

	// Update the validator delegations every 1 hour
	if _, err := scheduler.Every(1).Hour().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateValidatorsDelegations(stakingClient, cdc, db) })
	}); err != nil {
		return fmt.Errorf("error while setting up staking periodic operation: %s", err)
	}

	return nil
}

// updateValidatorsDelegations updates the current validators set and all their delegations, unbonding delegations
// and redelegations
func updateValidatorsDelegations(
	stakingClient stakingtypes.QueryClient, cdc codec.Marshaler, db *database.Db,
) error {
	height, err := db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting last block height: %s", err)
	}

	// Get the params
	params, err := db.GetStakingParams()
	if err != nil {
		return fmt.Errorf("error while getting staking params: %s", err)
	}

	validators, err := stakingutils.UpdateValidators(height, stakingClient, cdc, db)
	if err != nil {
		return fmt.Errorf("error while updating validators: %s", err)
	}

	// Update the delegations
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stakingutils.UpdateValidatorsDelegations(height, validators, stakingClient, db)
		wg.Done()
	}()

	// Update the unbonding delegations
	wg.Add(1)
	go func() {
		stakingutils.UpdateValidatorsUnbondingDelegations(height, params.BondDenom, validators, stakingClient, db)
		wg.Done()
	}()

	// Update the redelegations
	wg.Add(1)
	go func() {
		stakingutils.UpdateValidatorsRedelegations(height, params.BondDenom, validators, stakingClient, db)
		wg.Done()
	}()

	wg.Wait()
	return nil
}
