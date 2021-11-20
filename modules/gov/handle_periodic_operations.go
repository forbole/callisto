package gov

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v2/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "gov").Msg("setting up periodic tasks")

	// Update the community pool every 1 hour
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.handleMissingParams)
	}); err != nil {
		return fmt.Errorf("error while scheduling gov peridic operation: %s", err)
	}

	return nil
}

// handleMissingParams updates the params if any params table is empty
func (m *Module) handleMissingParams() error {
	lastBlock, err := m.db.GetLastBlock()
	if err != nil {
		return fmt.Errorf("error while getting last block height: %s", err)
	}
	lastHeight := lastBlock.Height

	// Update distribution params table if no value is stored
	existingDistrParams, err := m.db.GetDistributionParams()
	if err != nil {
		return fmt.Errorf("error while getting distribution params: %s", err)
	}
	if existingDistrParams == nil {
		err = m.distrModule.UpdateParams(lastHeight)
		if err != nil {
			return fmt.Errorf("error while updating distribution params: %s", err)
		}
	}

	// Update gov params table if no value is stored
	existingGovParams, err := m.db.GetGovParams()
	if err != nil {
		return fmt.Errorf("error while getting gov params: %s", err)
	}
	if existingGovParams == nil {
		err = m.UpdateParams(lastHeight)
		if err != nil {
			return fmt.Errorf("error while updating gov params: %s", err)
		}
	}

	// Update mint params table if no value is stored
	existingMintParams, err := m.db.GetMintParams()
	if err != nil {
		return fmt.Errorf("error while getting mint params: %s", err)
	}
	if existingMintParams == nil {
		err = m.mintModule.UpdateParams(lastHeight)
		if err != nil {
			return fmt.Errorf("error while updating mint params: %s", err)
		}
	}

	// Update slashing params table if no value is stored
	existingSlashingParams, err := m.db.GetSlashingParams()
	if err != nil {
		return fmt.Errorf("error while getting slashing params: %s", err)
	}
	if existingSlashingParams == nil {
		err = m.slashingModule.UpdateParams(lastHeight)
		if err != nil {
			return fmt.Errorf("error while updating slashing params: %s", err)
		}
	}

	// Update staking params table if no value is stored
	existingStakingParams, err := m.db.GetStakingParams()
	if err != nil {
		return fmt.Errorf("error while getting staking params: %s", err)
	}
	if existingStakingParams == nil {
		err = m.stakingModule.UpdateParams(lastHeight)
		if err != nil {
			return fmt.Errorf("error while updating staking params: %s", err)
		}
	}

	return nil
}
