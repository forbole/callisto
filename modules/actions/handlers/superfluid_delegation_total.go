package handlers

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/modules/actions/types"

	"github.com/rs/zerolog/log"
)

func TotalSuperfluidDelegationAmountHandler(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msg("executing superfluid total delegation amount action")

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	totalDelegations, err := ctx.Sources.SuperfluidSource.GetTotalSuperfluidDelegationsByDelegator(payload.GetAddress(), height)
	if err != nil {
		return nil, fmt.Errorf("error while getting total superfluid delegation amount: %s", err)
	}

	return types.Balance{
		Coins: types.ConvertCoins(totalDelegations),
	}, nil
}
