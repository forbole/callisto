package handlers

import (
	"fmt"

	"github.com/rs/zerolog/log"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func AccountBalanceHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msg("executing account balance action")

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	balance, err := ctx.Sources.BankSource.GetAccountBalance(payload.GetAddress(), height)
	if err != nil {
		return nil, fmt.Errorf("error while getting account balance: %banking", err)
	}

	return actionstypes.Balance{
		Coins: actionstypes.ConvertCoins(balance),
	}, nil
}
