package distribution

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v2/types"
)

// shouldUpdateDelegatorRewardsAmounts tells whether the delegators reward amounts should be updated at the given height
func (m *Module) shouldUpdateDelegatorRewardsAmounts(height int64) bool {
	interval := m.cfg.RewardsFrequency
	if interval == 0 {
		log.Debug().Str("module", "distribution").Msg("delegator rewards refresh interval set to 0. Skipping refresh")
		return false
	}

	hasRewards, err := m.db.HasDelegatorRewards()
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).Int64("height", height).
			Msg("error while checking delegators reward")
		return false
	}

	return !hasRewards || height%interval == 0
}

// refreshDelegatorsRewardsAmounts refreshes the rewards associated with all the delegators for the given height,
// deleting the ones existing and downloading them from scratch.
func (m *Module) refreshDelegatorsRewardsAmounts(height int64) {
	// Get the delegators
	delegators, err := m.db.GetDelegators()
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).Int64("height", height).
			Msg("error while getting delegators")
	}

	if len(delegators) == 0 {
		log.Debug().Str("module", "distribution").Int64("height", height).
			Msg("no delegations found, make sure you are calling this module after the staking module")
		return
	}

	// Get the rewards
	for _, delegator := range delegators {
		err = m.RefreshDelegatorRewards(height, delegator)
		if err != nil {
			log.Error().Str("module", "distribution").Err(err).Int64("height", height).
				Str("delegator", delegator).Msg("error while updating delegator rewards")
		}
	}
}

// RefreshDelegatorRewards refreshes the rewards associated to the given delegator for the given height,
// deleting the ones existing and downloading them from scratch.
func (m *Module) RefreshDelegatorRewards(height int64, delegator string) error {
	rewards, err := m.getDelegatorRewardsAmounts(height, delegator)
	if err != nil {
		return fmt.Errorf("error while refreshing delegator rewards: %s", err)
	}

	err = m.db.SaveDelegatorsRewardsAmounts(height, delegator, rewards)
	if err != nil {
		return fmt.Errorf("error while saving delegators rewards amounts: %s", err)
	}

	return nil
}

// getDelegatorRewardsAmounts returns the rewards for the given delegator at the given height
func (m *Module) getDelegatorRewardsAmounts(height int64, delegator string) ([]types.DelegatorRewardAmount, error) {
	rows, err := m.source.DelegatorTotalRewards(delegator, height)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator reward: %s", err)
	}

	withdrawAddr, err := m.source.DelegatorWithdrawAddress(delegator, height)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator rewards: %s", err)
	}

	var rewards = make([]types.DelegatorRewardAmount, len(rows))
	for index, reward := range rows {
		rewards[index] = types.NewDelegatorRewardAmount(reward.ValidatorAddress, withdrawAddr, reward.Reward)
	}

	return rewards, nil
}
