package bank

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v5/modules/utils"
)

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "bank").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(10).Minutes().Do(func() {
		utils.WatchMethod(m.UpdateSupply)
	}); err != nil {
		return fmt.Errorf("error while setting up bank periodic operation: %s", err)
	}

	return nil
}

// UpdateSupply updates the supply of all the tokens
func (m *Module) UpdateSupply() error {
	log.Trace().Str("module", "bank").Str("operation", "total supply").
		Msg("updating total supply")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	supply, err := m.keeper.GetSupply(height)
	if err != nil {
		return err
	}

	return m.db.SaveSupply(supply, height)
}
