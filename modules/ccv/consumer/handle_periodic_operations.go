package consumer

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "ccvconsumer").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(15).Minutes().Do(func() {
		utils.WatchMethod(m.updateNextFeeDistribution)
	}); err != nil {
		return fmt.Errorf("error while setting up ccvconsumer periodic operation: %s", err)
	}

	return nil
}

// updateNextFeeDistribution updates next fee distribution estimate in database
func (m *Module) updateNextFeeDistribution() error {
	block, err := m.db.GetLastBlock()
	if err != nil {
		return fmt.Errorf("error while getting last block: %s", err)
	}

	log.Trace().Str("module", "ccvconsumer").Int64("height", block.Height).
		Msg("updating next fee distribution")

	nextFeeDistributionEstimate, err := m.source.GetNextFeeDistribution(block.Height)
	if err != nil {
		return fmt.Errorf("error while getting next fee distribution estimate: %s", err)
	}

	return m.db.SaveNextFeeDistributionEstimate(block.Height, nextFeeDistributionEstimate)
}
