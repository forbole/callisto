package consumer

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "ccvconsumer").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.UpdateCCVValidators)
	}); err != nil {
		return fmt.Errorf("error while setting up ccvconsumer period operations: %s", err)
	}

	return nil
}

// UpdateCCVValidators updated ccv validators details inside the database
// to create relationship between consumer and provider chain
func (m *Module) UpdateCCVValidators() error {
	log.Debug().
		Str("module", "ccvconsumer").
		Msg("updating ccv validators details")

	block, err := m.db.GetLastBlock()
	if err != nil {
		return fmt.Errorf("error while getting last block: %s", err)
	}

	err = m.UpdateCcvValidators(block.Height)
	if err != nil {
		return err
	}

	return nil

}