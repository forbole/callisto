package handlers

import (
	"fmt"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/rs/zerolog/log"
)

func ValidatorCommissionAmountHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msg("executing validator commission action")

	// Get latest node height
	height, err := ctx.GetHeight(nil)
	if err != nil {
		return nil, err
	}

	// Get validator total commission value
	commission, err := ctx.Sources.DistrSource.ValidatorCommission(payload.GetAddress(), height)
	if err != nil {
		return nil, fmt.Errorf("error while getting validator commission: %s", err)
	}

	return actionstypes.ValidatorCommissionAmount{
		Coins: actionstypes.ConvertDecCoins(commission),
	}, nil
}
