package handlers

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/modules/actions/types"

	"github.com/rs/zerolog/log"
)

func SuperfluidDelegationsHandler(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msg("executing superfluid delegations action")

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	delegations, err := ctx.Sources.SuperfluidSource.GetSuperfluidDelegationsByDelegator(payload.GetAddress(), height)
	if err != nil {
		return nil, fmt.Errorf("error while getting account superfluid delegations: %s", err)
	}

	return delegations, nil
}
