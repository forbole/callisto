package oracle

import (
	"github.com/forbole/bdjuno/v2/modules/utils"
	"github.com/forbole/bdjuno/v2/types"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "oracle").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").StartImmediately().Do(func() {
		utils.WatchMethod(m.updateOracleParams)
	}); err != nil {
		return err
	}

	return nil
}

// updateOracleParams fetches from the REST APIs the latest value for
// oracle params, and saves it inside the database.
func (m *Module) updateOracleParams() error {
	log.Debug().Str("module", "oracle").Msg("getting oracle params data")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the params
	params, err := m.source.GetParams(height)
	if err != nil {
		return err
	}

	return m.db.SaveOracleParams(types.NewOracleParams(params, height), height)
}
