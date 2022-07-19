package mint

import (
	"fmt"
	"strconv"
	"time"

	"github.com/forbole/bdjuno/v3/modules/utils"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	creminttypes "github.com/crescent-network/crescent/v2/x/mint/types"
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

	annualProvision := getCurrentAnnualProvision(block.Timestamp, mintParams.InflationSchedules)
	if annualProvision == 0 {
		return nil
	}

	// Get current total supply of uCRE
	supply, err := m.db.GetSupply("ucre")
	if err != nil {
		return err
	}

	// Convert supply to int64
	supplyInt, err := strconv.ParseInt(supply.Amount, 10, 64)
	if err != nil {
		return fmt.Errorf("error while converting supply string to int64: %s", err)
	}

	// Calculate the inflation: annual provision / current total supply
	inflation := float64(annualProvision) / float64(supplyInt)

	return m.db.SaveInflation(inflation, block.Height)
}

// getCurrentAnnualProvision gets the new annual provision if block time enters into new inflation schedule, and if not, returns 0
func getCurrentAnnualProvision(
	blockTime time.Time, inflationSchedules []creminttypes.InflationSchedule,
) int64 {
	for _, schedule := range inflationSchedules {
		if blockTime.Year() == schedule.StartTime.Year() &&
			blockTime.Month() == schedule.StartTime.Month() &&
			blockTime.Day() == schedule.StartTime.Day() {

			return schedule.Amount.Int64()
		}
	}

	return 0
}
