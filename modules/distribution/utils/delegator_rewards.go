package utils

import (
	"context"
	"fmt"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types/config"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/types"
)

// UpdateDelegatorsRewardsAmounts updates the delegators commission amounts
func UpdateDelegatorsRewardsAmounts(cfg *config.Config, height int64, client distrtypes.QueryClient, db *database.Db) {
	interval := cfg.GetDistributionConfig().GetDistributionFrequency()
	if interval == 0 {
		log.Debug().Str("module", "distribution").Msg("delegator rewards refresh interval set to 0. Skipping refresh")
		return
	}

	hasRewards, error := db.HasDelegatorRewards()
	if error != nil {
		log.Error().Str("module", "distribution").Err(error).Int64("height", height).
			Msg("error while checking delegators reward")
	}

	if !hasRewards || height%interval == 0 {
		go updateDelegatorsRewards(height, client, db)
	}

}

func updateDelegatorsRewards(height int64, client distrtypes.QueryClient, db *database.Db) {
	log.Debug().Str("module", "distribution").Int64("height", height).
		Msg("updating delegators rewards")

	// Get the delegators
	delegators, err := db.GetDelegators()
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
		go updateDelegatorCommission(height, delegator, client, db)
	}
}

func updateDelegatorCommission(height int64, delegator string, distrClient distrtypes.QueryClient, db *database.Db) {
	rewards, err := GetDelegatorRewards(height, delegator, distrClient)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).Str("delegator", delegator).
			Msg("error while getting delegator rewards")
	}

	err = db.SaveDelegatorsRewardsAmounts(rewards)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).Str("delegator", delegator).
			Msg("error while saving delegator rewards")
	}
}

// GetDelegatorRewards returns the current rewards for the given delegator
func GetDelegatorRewards(height int64, delegator string, distrClient distrtypes.QueryClient) ([]types.DelegatorRewardAmount, error) {
	header := client.GetHeightRequestHeader(height)

	rewardsRes, err := distrClient.DelegationTotalRewards(
		context.Background(),
		&distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegator},
		header,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator reward: %s", err)
	}

	withdrawAddressRes, err := distrClient.DelegatorWithdrawAddress(
		context.Background(),
		&distrtypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: delegator},
		header,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator rewards: %s", err)
	}

	var rewards = make([]types.DelegatorRewardAmount, len(rewardsRes.Rewards))
	for index, reward := range rewardsRes.Rewards {
		rewards[index] = types.NewDelegatorRewardAmount(
			delegator,
			reward.ValidatorAddress,
			withdrawAddressRes.WithdrawAddress,
			reward.Reward,
			height,
		)
	}
	return rewards, nil
}

// RefreshDelegatorRewards refreshes the rewards associated to the given delegator for the given height,
// deleting the ones existing and downloading them from scratch.
func RefreshDelegatorRewards(height int64, delegator string, distrClient distrtypes.QueryClient, db *database.Db) error {
	rewards, err := GetDelegatorRewards(height, delegator, distrClient)
	if err != nil {
		return fmt.Errorf("error while refreshing delegator rewards: %s", err)
	}

	err = db.DeleteDelegatorRewardsAmount(delegator, height)
	if err != nil {
		return fmt.Errorf("error deleting the delegator rewards amount: %s", err)
	}

	err = db.SaveDelegatorsRewardsAmounts(rewards)
	if err != nil {
		return fmt.Errorf("error while saving delegators rewards amounts: %s", err)
	}

	return nil
}
