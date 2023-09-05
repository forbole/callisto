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
		utils.WatchMethod(m.updateProposalStakingPoolSnapshots)
	}); err != nil {
		return fmt.Errorf("error while setting up gov period operations: %s", err)
	}

	// refresh proposal validators status snapshots every 5 mins
	if _, err := scheduler.Every(5).Minutes().Do(func() {
		utils.WatchMethod(m.updateProposalValidatorsStatusSnapshot)
	}); err != nil {
		return fmt.Errorf("error while setting up gov period operations: %s", err)
	}

	// refresh proposal tally results every 5 mins
	if _, err := scheduler.Every(5).Minutes().Do(func() {
		utils.WatchMethod(m.updateProposalTallyResult)
	}); err != nil {
		return fmt.Errorf("error while setting up gov period operations: %s", err)
	}

	return nil
}

// updateProposalStakingPoolSnapshots updated staing pool snapshot for
// every active proposal
func (m *Module) updateProposalStakingPoolSnapshots() error {
	log.Debug().Str("module", "gov").Msg("refreshing proposal staking pool snapshots")

	blockTime, err := m.db.GetLastBlockTimestamp()
	if err != nil {
		return err
	}

	ids, err := m.db.GetOpenProposalsIds(blockTime)
	if err != nil {
		log.Error().Err(err).Str("module", "gov").Msg("error while getting open proposals ids")
	}

	// Get the latest block height from db
	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height from db: %s", err)
	}

	for _, proposalID := range ids {
		err = m.UpdateProposalStakingPoolSnapshot(height, proposalID)
		if err != nil {
			return fmt.Errorf("error while updating proposal %d staking pool snapshots: %s", proposalID, err)
		}
	}

	return nil
}

// updateProposalTallyResult updates the tally for active proposals
func (m *Module) updateProposalTallyResult() error {
	log.Debug().Str("module", "gov").Msg("refreshing proposal tally results")
	blockTime, err := m.db.GetLastBlockTimestamp()
	if err != nil {
		return err
	}

	ids, err := m.db.GetOpenProposalsIds(blockTime)
	if err != nil {
		log.Error().Err(err).Str("module", "gov").Msg("error while getting open proposals ids")
	}

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	for _, proposalID := range ids {
		err = m.UpdateProposalTallyResult(proposalID, height)
		if err != nil {
			return fmt.Errorf("error while updating proposal %d tally result : %s", proposalID, err)
		}
	}

	return nil
}
