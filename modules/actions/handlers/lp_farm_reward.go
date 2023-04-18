package handlers

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/modules/actions/types"

	"github.com/rs/zerolog/log"
)

func LPFarmRewardHandler(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msg("executing lp farm rewards action")

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get delegator's lp farm total rewards
	rewards, err := ctx.Sources.LPFarmSource.GetTotalLPFarmRewards(payload.GetAddress(), height)
	if err != nil {
		return nil, fmt.Errorf("error while getting farmer total rewards: %s", err)
	}

	farmerRewards := make([]types.LPFarmReward, len(rewards))
	for index, reward := range rewards {
		farmerRewards[index] = types.LPFarmReward{
			Coins:         []types.Coin{{Amount: reward.Amount.String(), Denom: reward.Denom}},
			FarmerAddress: payload.GetAddress(),
		}
	}

	return farmerRewards, nil
}
