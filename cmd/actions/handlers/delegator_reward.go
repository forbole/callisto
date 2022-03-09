package handlers

import (
	"fmt"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func DelegationRewardHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get delegator's total rewards
	rewards, err := ctx.Sources.DistrSource.DelegatorTotalRewards(payload.GetAddress(), height)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator total rewards: %s", err)
	}

	delegationRewards := make([]actionstypes.DelegationReward, len(rewards))
	for index, rew := range rewards {
		delegationRewards[index] = actionstypes.DelegationReward{
			Coins:            actionstypes.ConvertDecCoins(rew.Reward),
			ValidatorAddress: rew.ValidatorAddress,
		}
	}

	return delegationRewards, nil
}
