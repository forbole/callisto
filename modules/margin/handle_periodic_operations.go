package margin

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/modules/utils"
	"github.com/forbole/bdjuno/v3/types"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "margin").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.updateParams)
	}); err != nil {
		return err
	}

	return nil
}

// updateParams fetches the latest x/margin module params
// and saves it inside the database.
func (m *Module) updateParams() error {
	log.Debug().
		Str("module", "margin").
		Str("operation", "params").
		Msg("updating x/margin params")

	block, err := m.db.GetLastBlock()
	if err != nil {
		return fmt.Errorf("error while getting last block: %s", err)
	}

	// Get the latest params
	params, err := m.source.GetParams(block.Height)
	if err != nil {
		return err
	}

	return m.db.SaveMarginParams(types.NewMarginParams(params, block.Height))
}
