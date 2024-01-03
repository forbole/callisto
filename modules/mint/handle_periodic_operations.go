package mint

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/forbole/bdjuno/v4/modules/utils"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "mint").Msg("setting up periodic tasks")

	// Fetch inflation once per day at midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.UpdateInflation)
	}); err != nil {
		return fmt.Errorf("error while setting up inflation periodic operations: %s", err)
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

	block, err := m.db.GetLastBlockHeightAndTimestamp()
	if err != nil {
		return err
	}

	mintParams, err := m.source.Params(block.Height)
	if err != nil {
		return err
	}

	epochProvisions, err := m.source.EpochProvisions(block.Height)
	if err != nil {
		return err
	}

	epochProvisionsFloat := new(big.Float).SetInt64(epochProvisions.RoundInt64())
	epochPeriod := new(big.Float).SetInt64(mintParams.ReductionPeriodInEpochs)

	totalSupply, err := m.db.GetTotalSupply()
	if err != nil {
		return err
	}
	supply, err := strconv.ParseInt(totalSupply, 10, 64)
	if err != nil {
		return err
	}

	totalProvisions := new(big.Float).Mul(epochProvisionsFloat, epochPeriod)
	inflation := new(big.Float).Quo(totalProvisions, new(big.Float).SetInt64(supply))

	return m.db.SaveInflation(inflation.String(), block.Height)

}
