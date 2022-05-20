package mint

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/forbole/bdjuno/v2/modules/utils"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "mint").Msg("setting up periodic tasks")

	// Fetch inflation once per day at midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.updateInflation)
	}); err != nil {
		return fmt.Errorf("error while setting up inflation periodic operations: %s", err)
	}

	return nil
}

// updateInflation caluclates current inflation value and stores it in database
func (m *Module) updateInflation() error {
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

	epochProvisions := new(big.Float).SetInt64(mintParams.GenesisEpochProvisions.RoundInt64())
	epochPeriod := new(big.Float).SetInt64(mintParams.ReductionPeriodInEpochs)

	totalSupply, err := m.db.GetTotalSupply()
	if err != nil {
		return err
	}
	supply, err := strconv.ParseInt(totalSupply, 10, 64)
	if err != nil {
		return err
	}

	totalProvisions := new(big.Float).Mul(epochProvisions, epochPeriod)
	inflation := new(big.Float).Quo(totalProvisions, new(big.Float).SetInt64(supply))

	return m.db.SaveInflation(inflation.String(), height)

}
