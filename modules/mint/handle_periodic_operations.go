package mint

import (
	"fmt"
	"strconv"

	"github.com/forbole/bdjuno/v3/modules/utils"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "mint").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.updateInflation)
	}); err != nil {
		return err
	}

	return nil
}

// updateInflation fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func (m *Module) updateInflation() error {
	log.Debug().
		Str("module", "mint").
		Str("operation", "inflation").
		Msg("getting inflation data")

	block, err := m.db.GetLastBlock()
	if err != nil {
		return err
	}

	// Get mint params
	mintParams, err := m.db.GetMintParams()
	if err != nil {
		return err
	}

	// Get the new annual provision if entering into new inflation schedule
	var annualProvision int64
	for _, schedule := range mintParams.InflationSchedules {
		if block.Timestamp.Year() == schedule.StartTime.Year() &&
			block.Timestamp.Month() == schedule.StartTime.Month() &&
			block.Timestamp.Day() == schedule.StartTime.Day() {

			annualProvision = schedule.Amount.Int64()
		} else {

			return nil
		}
	}

	// Get current total supply of uCRE
	currTotalSupply, err := m.db.GetSupply("ucre")
	if err != nil {
		return err
	}

	// Convert supply to int64
	supplyInt, err := strconv.ParseInt(currTotalSupply.Amount, 10, 64)
	if err != nil {
		return fmt.Errorf("error while converting supply to int64: %s", err)
	}

	// Calculate the inflation: annual provision / current total supply
	inflation := float64(annualProvision) / float64(supplyInt)

	return m.db.SaveInflation(inflation, block.Height)
}
