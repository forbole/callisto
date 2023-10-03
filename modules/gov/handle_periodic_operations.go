package gov

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "gov").Msg("setting up periodic tasks")

	// refresh proposal staking pool snapshots every 5 mins
	// (set the same interval as staking pool periodic ops)
	if _, err := scheduler.Every(5).Minutes().Do(func() {
		utils.WatchMethod(m.UpdateProposalsStakingPoolSnapshot)
	}); err != nil {
		return fmt.Errorf("error while setting up gov period operations: %s", err)
	}

	// refresh proposal tally results every 5 mins
	if _, err := scheduler.Every(5).Minutes().Do(func() {
		utils.WatchMethod(m.UpdateProposalsTallyResults)
	}); err != nil {
		return fmt.Errorf("error while setting up gov period operations: %s", err)
	}

	return nil
}
