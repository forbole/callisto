package handlers

import (
	"fmt"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func DelegatorWithdrawAddressHandler(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	// Get latest node height
	height, err := ctx.GetHeight(nil)
	if err != nil {
		return nil, err
	}

	// Get delegator's total rewards
	withdrawAddress, err := ctx.Sources.DistrSource.DelegatorWithdrawAddress(payload.GetAddress(), height)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator withdraw address: %s", err)
	}

	return actionstypes.Address{
		Address: withdrawAddress,
	}, nil
}
