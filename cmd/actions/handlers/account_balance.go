package handlers

import (
	"fmt"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/rs/zerolog/log"
)

func AccountBalanceHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	log.Debug().Str("action", "account balance").
		Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height)

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	balance, err := ctx.Sources.BankSource.GetAccountBalance(payload.GetAddress(), height)
	if err != nil {
		return nil, fmt.Errorf("error while getting account balance: %s", err)
	}

	return actionstypes.Balance{
		Coins: actionstypes.ConvertCoins(balance),
	}, nil
}
