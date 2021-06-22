package utils

import (
	"context"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"
)

// UpdateDelegatorsRewardsAmounts updates the delegators commission amounts
func UpdateDelegatorsRewardsAmounts(height int64, client distrtypes.QueryClient, db *database.Db) {
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
		go getDelegatorCommission(height, delegator, client, db)
	}
}

func getDelegatorCommission(height int64, delegator string, distrClient distrtypes.QueryClient, db *database.Db) {
	header := client.GetHeightRequestHeader(height)

	rewardsRes, err := distrClient.DelegationTotalRewards(
		context.Background(),
		&distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegator},
		header,
	)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).Str("delegator", delegator).
			Msg("error while getting delegator reward")
		return
	}

	withdrawAddressRes, err := distrClient.DelegatorWithdrawAddress(
		context.Background(),
		&distrtypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: delegator},
		header,
	)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).Str("delegator", delegator).
			Msg("error while getting delegator withdraw address")
		return
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

	err = db.SaveDelegatorsRewardsAmounts(rewards)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).Str("delegator", delegator).
			Msg("error while saving delegator rewards")
	}
}
