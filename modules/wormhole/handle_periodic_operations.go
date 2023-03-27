package wormhole

import (
	"github.com/forbole/bdjuno/v4/modules/utils"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "wormhole").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight to update guardian
	// set in database
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.UpdateGuardianSet)
	}); err != nil {
		return err
	}

	// Setup a cron job to run every midnight to update guardian
	// validators list in database
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.UpdateGuardianValidators)
	}); err != nil {
		return err
	}
	return nil
}

// UpdateGuardianSet fetches the latest
// guardian set, and saves it inside the database.
func (m *Module) UpdateGuardianSet() error {
	log.Debug().
		Str("module", "wormhole").
		Msg("updating wormhole guardian set")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the latest guardian set
	guardianSetList, err := m.source.GetGuardianSetAll(height)
	if err != nil {
		return err
	}

	return m.db.SaveGuardianSetList(guardianSetList, height)
}

// UpdateGuardianValidators fetches the latest
// guardian validators, and saves them inside the database.
func (m *Module) UpdateGuardianValidators() error {
	log.Debug().
		Str("module", "wormhole").
		Msg("updating wormhole guardian validators")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the latest guardian validators list
	guardianValidatorsList, err := m.source.GetGuardianValidatorAll(height)
	if err != nil {
		return err
	}

	return m.db.SaveGuardianValidatorList(guardianValidatorsList, height)
}
