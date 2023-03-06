package bank

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "bank").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(10).Minutes().Do(func() {
		utils.WatchMethod(m.updateSupply)
	}); err != nil {
		return fmt.Errorf("error while setting up bank periodic operation: %s", err)
	}

	return nil
}

// updateSupply updates the supply of all the tokens
func (m *Module) updateSupply() error {
	log.Trace().Str("module", "bank").Str("operation", "total supply").
		Msg("updating total supply")

	block, err := m.db.GetLastBlock()
	if err != nil {
		return fmt.Errorf("error while getting last block: %s", err)
	}

	params, err := m.stakingModule.GetParams(block.Height)
	if err != nil {
		return fmt.Errorf("error while getting bond denom type: %s", err)
	}

	supply, err := m.keeper.GetSupply(block.Height, params.BondDenom)
	if err != nil {
		return err
	}

	return m.db.SaveSupply(supply, block.Height)
}
