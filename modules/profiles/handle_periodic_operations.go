package profiles

import (
	"github.com/forbole/bdjuno/v2/modules/utils"
	"github.com/forbole/bdjuno/v2/types"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "profiles").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.updateParam)
	}); err != nil {
		return err
	}

	return nil
}

// updateInflation fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func (m *Module) updateParam() error {
	log.Debug().Str("module", "profiles").Str("operation", "profiles").
		Msg("getting profiles params")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the inflation
	inflation, err := m.source.GetParams(height)
	if err != nil {
		return err
	}

	return m.db.SaveProfilesParams(types.NewProfilesParams(inflation, height))
}
