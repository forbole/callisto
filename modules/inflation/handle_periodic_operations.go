package mint

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/modules/utils"
	"github.com/forbole/bdjuno/v3/types"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "inflation").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.updateInflationData)
	}); err != nil {
		return err
	}

	return nil
}

// updateInflation fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func (m *Module) updateInflationData() error {
	log.Debug().
		Str("module", "inflation").
		Msg("getting inflation data")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	evmosInfationData, err := m.getInflationData(height)
	if err != nil {
		return fmt.Errorf("error while storing evmos inflation data: %s", err)
	}

	return m.db.SaveEvmosInflationData(evmosInfationData, height)
}

func (m *Module) getInflationData(height int64) (*types.EvmosInflationData, error) {
	circulatingSupply, err := m.source.CirculatingSupply(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting inflation rate: %s", err)
	}

	epochMintProvision, err := m.source.EpochMintProvision(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting epoch mint provision: %s", err)
	}

	inflationRate, err := m.source.InflationRate(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting inflation rate: %s", err)
	}

	inflationPeriod, err := m.source.InflationPeriod(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting inflation period: %s", err)
	}

	skippedEpochs, err := m.source.SkippedEpochs(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting skipped epochs: %s", err)
	}

	return types.NewEvmosInflationData(
		circulatingSupply, epochMintProvision, inflationRate, inflationPeriod, skippedEpochs, height,
	), nil
}
