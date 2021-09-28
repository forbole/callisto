package distribution

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/types"
)

// refreshDelegatorsRewardsAmounts refreshes the rewards associated with all the delegators for the given height,
// deleting the ones existing and downloading them from scratch.
func (m *Module) refreshDelegatorsRewardsAmounts(height int64) {
	interval := m.cfg.DistributionFrequency
	if interval == 0 {
		log.Debug().Str("module", "distribution").Msg("delegator rewards refresh interval set to 0. Skipping refresh")
		return
	}

	hasRewards, err := m.db.HasDelegatorRewards()
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).Int64("height", height).
			Msg("error while checking delegators reward")
	}

	if hasRewards && height%interval != 0 {
		return
	}

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

		// Refresh the delegators using a goroutine to improve efficiency
		go func(delegator string) {
			err = m.refreshDelegatorRewards(height, delegator)
			if err != nil {
				log.Error().Str("module", "distribution").Err(err).Int64("height", height).
					Str("delegator", delegator).Msg("error while updating delegator rewards")
			}
		}(delegator)
	}
}

// refreshDelegatorRewards refreshes the rewards associated to the given delegator for the given height,
// deleting the ones existing and downloading them from scratch.
func (m *Module) refreshDelegatorRewards(height int64, delegator string) error {
	rewards, err := m.getDelegatorRewardsAmounts(height, delegator)
	if err != nil {
		return fmt.Errorf("error while refreshing delegator rewards: %s", err)
	}

	err = m.db.DeleteDelegatorRewardsAmount(delegator, height)
	if err != nil {
		return fmt.Errorf("error deleting the delegator rewards amount: %s", err)
	}

	err = m.db.SaveDelegatorsRewardsAmounts(rewards)
	if err != nil {
		return fmt.Errorf("error while saving delegators rewards amounts: %s", err)
	}

	return nil
}

// getDelegatorRewardsAmounts returns the rewards for the given delegator at the given height
func (m *Module) getDelegatorRewardsAmounts(height int64, delegator string) ([]types.DelegatorRewardAmount, error) {
	rews, err := m.source.DelegatorTotalRewards(delegator, height)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator reward: %s", err)
	}

	withdrawAddr, err := m.source.DelegatorWithdrawAddress(delegator, height)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator rewards: %s", err)
	}

	var rewards = make([]types.DelegatorRewardAmount, len(rews))
	for index, reward := range rews {
		rewards[index] = types.NewDelegatorRewardAmount(
			delegator,
			reward.ValidatorAddress,
			withdrawAddr,
			reward.Reward,
			height,
		)
	}

	return rewards, nil
}
