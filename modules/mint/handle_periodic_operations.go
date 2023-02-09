package mint

import (
	"math/big"
	"strconv"

	"github.com/forbole/bdjuno/v4/modules/utils"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "mint").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.UpdateInflation)
	}); err != nil {
		return err
	}

	return nil
}

// updateInflation fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func (m *Module) UpdateInflation() error {
	log.Debug().
		Str("module", "mint").
		Str("operation", "inflation").
		Msg("getting inflation data")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	mintParams, err := m.source.Params(height)
	if err != nil {
		return err
	}

	// Get the epoch provisions
	epochProvisions, err := m.source.GetEpochProvisions(height)
	if err != nil {
		return err
	}

	epochProvisionsFloat := new(big.Float).SetInt64(epochProvisions.RoundInt64())

	// Normally 365 days
	reductionPeriod := new(big.Float).SetInt64(mintParams.ReductionPeriodInEpochs)

	// Get total supply
	totalSupply, err := m.db.GetTotalSupply()
	if err != nil {
		return err
	}
	supply, err := strconv.ParseInt(totalSupply, 10, 64)
	if err != nil {
		return err
	}

	// Provision per epoch * reduction period(365 days) = yearly provision
	totalProvisions := new(big.Float).Mul(epochProvisionsFloat, reductionPeriod)

	// Inflation = yearly provision / current total supply
	inflation := new(big.Float).Quo(totalProvisions, new(big.Float).SetInt64(supply))

	return m.db.SaveInflation(inflation.String(), height)
}
