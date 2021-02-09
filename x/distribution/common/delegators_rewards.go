package common

import (
	"context"
	"sync"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	bdistrtypes "github.com/forbole/bdjuno/x/distribution/types"
	"github.com/forbole/bdjuno/x/utils"
)

// UpdateDelegatorsRewardsAmounts updates the delegators commission amounts
func UpdateDelegatorsRewardsAmounts(height int64, client distrtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "distribution").Int64("height", height).
		Msg("updating delegators rewards")

	// Get the delegators
	delegators, err := db.GetDelegatorsForHeight(height)
	if err != nil {
		return err
	}

	if len(delegators) == 0 {
		log.Debug().Str("module", "distribution").Int64("height", height).
			Msg("no delegations found, make sure you are calling this module after the staking module")
		return nil
	}

	// Get the rewards
	var wg sync.WaitGroup
	var out = make(chan bdistrtypes.DelegatorRewardAmount)
	for _, delegator := range delegators {
		wg.Add(1)
		go getDelegatorCommission(height, delegator, client, db, out, &wg)
	}

	// We need to call wg.Wait inside another goroutine in order to solve the hanging bug that's described here:
	// https://dev.to/sophiedebenedetto/synchronizing-go-routines-with-channels-and-waitgroups-3ke2
	go func() {
		wg.Wait()
		close(out)
	}()

	var rewards []bdistrtypes.DelegatorRewardAmount
	for com := range out {
		rewards = append(rewards, com)
	}

	// Save the rewards
	return db.SaveDelegatorsRewardsAmounts(rewards, height)
}

func getDelegatorCommission(
	height int64,
	delegator string,
	client distrtypes.QueryClient,
	db *database.BigDipperDb,
	out chan<- bdistrtypes.DelegatorRewardAmount,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	res, err := client.DelegationTotalRewards(
		context.Background(),
		&distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegator},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).Str("delegator", delegator).
			Msg("error while getting delegator reward")
		return
	}

	for _, reward := range res.Rewards {
		consAddr, err := db.GetValidatorConsensusAddress(reward.ValidatorAddress)
		if err != nil {
			log.Error().Str("module", "distribution").Err(err).
				Int64("height", height).Str("delegator", delegator).
				Msg("error while getting delegator reward")
			return
		}

		// Send the reward amount back
		out <- bdistrtypes.NewDelegatorRewardAmount(consAddr.String(), delegator, reward.Reward)
	}
}
