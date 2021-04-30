package common

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	bdistrtypes "github.com/forbole/bdjuno/x/distribution/types"
	"github.com/forbole/bdjuno/x/utils"
)

// UpdateDelegatorsRewardsAmounts updates the delegators commission amounts
func UpdateDelegatorsRewardsAmounts(height int64, client distrtypes.QueryClient, db *database.BigDipperDb) {
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

func getDelegatorCommission(
	height int64, delegator string, client distrtypes.QueryClient, db *database.BigDipperDb,
) {
	header := utils.GetHeightRequestHeader(height)

	rewardsRes, err := client.DelegationTotalRewards(
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

	withdrawAddressRes, err := client.DelegatorWithdrawAddress(
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

	var rewards = make([]bdistrtypes.DelegatorRewardAmount, len(rewardsRes.Rewards))
	for index, reward := range rewardsRes.Rewards {
		consAddr, err := db.GetValidatorConsensusAddress(reward.ValidatorAddress)
		if err != nil {
			log.Error().Str("module", "distribution").Err(err).
				Int64("height", height).Str("delegator", delegator).
				Msg("error while getting delegator reward")
			return
		}

		// Send the reward amount back
		rewards[index] = bdistrtypes.NewDelegatorRewardAmount(
			delegator,
			consAddr.String(),
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
