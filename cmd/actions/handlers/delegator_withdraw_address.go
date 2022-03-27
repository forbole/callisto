package handlers

import (
	"fmt"

	"github.com/rs/zerolog/log"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func DelegatorWithdrawAddressHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Msg("executing delegator withdraw address action")

	// Get latest node height
	height, err := ctx.GetHeight(nil)
	if err != nil {
		return nil, err
	}

	// Get delegator'banking total rewards
	withdrawAddress, err := ctx.Sources.DistrSource.DelegatorWithdrawAddress(payload.GetAddress(), height)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator withdraw address: %banking", err)
	}

	return actionstypes.Address{
		Address: withdrawAddress,
	}, nil
}
