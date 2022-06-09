package shield

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v3/modules/utils"
	"github.com/forbole/bdjuno/v3/types"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "shield").Msg("setting up periodic tasks")

	// Fetch the pool providers every 10 mins
	if _, err := scheduler.Every(10).Minutes().Do(func() {
		utils.WatchMethod(m.updatePoolProviders)
	}); err != nil {
		return fmt.Errorf("error while setting up shield pool providers period operations: %s", err)
	}

	// Fetch the shield status every 10 mins
	if _, err := scheduler.Every(10).Minutes().Do(func() {
		utils.WatchMethod(m.updateShieldStatus)
	}); err != nil {
		return fmt.Errorf("error while setting up shield status period operations: %s", err)
	}
	return nil
}

// updatePoolProviders allows to get the most up-to-date pool providers
func (m *Module) updatePoolProviders() error {

	block, err := m.db.GetLastBlock()
	if err != nil {
		return fmt.Errorf("error while getting last block: %s", err)
	}

	providers, err := m.source.GetPoolProviders(block.Height)
	if err != nil {
		return err
	}

	for _, provider := range providers {
		err := m.db.SaveShieldProvider(types.NewShieldProvider(provider.Address, provider.Collateral.Int64(),
			provider.DelegationBonded.Int64(), provider.Rewards.Native, provider.Rewards.Foreign,
			provider.TotalLocked.Int64(), provider.Withdrawing.Int64(), block.Height))
		if err != nil {
			return err
		}
	}
	return nil
}

// updateShieldStatus allows to get the most up-to-date shield status
func (m *Module) updateShieldStatus() error {

	block, err := m.db.GetLastBlock()
	if err != nil {
		return fmt.Errorf("error while getting last block: %s", err)
	}

	shieldStatus, err := m.source.GetShieldStatus(block.Height)
	if err != nil {
		return err
	}

	err = m.db.SaveShieldStatus(types.NewShieldStatus(shieldStatus.GlobalShieldStakingPool,
		shieldStatus.CurrentServiceFees.Native, shieldStatus.CurrentServiceFees.Foreign,
		shieldStatus.RemainingServiceFees.Native, shieldStatus.RemainingServiceFees.Foreign,
		shieldStatus.TotalCollateral, shieldStatus.TotalShield, shieldStatus.TotalWithdrawing, block.Height))
	if err != nil {
		return err
	}

	return nil
}
