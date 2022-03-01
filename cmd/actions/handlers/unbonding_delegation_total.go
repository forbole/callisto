package handlers

import (
	"fmt"
	"math/big"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func UnbondingDelegationsTotal(ctx *actionstypes.Context, payload *actionstypes.Payload) (interface{}, error) {
	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get all unbonding delegations for given delegator address
	unbondingDelegations, err := ctx.Sources.StakingSource.GetUnbondingDelegations(height, payload.GetAddress(), nil)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator delegations: %s", err)
	}

	// Get the bond denom type
	params, err := ctx.Sources.StakingSource.GetParams(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting bond denom type: %s", err)
	}

	// Add up total value of unbonding delegations
	var totalAmount = big.NewInt(0)
	for _, eachUnbondingDelegation := range unbondingDelegations.UnbondingResponses {
		for _, entry := range eachUnbondingDelegation.Entries {
			totalAmount = totalAmount.Add(totalAmount, entry.Balance.BigInt())
		}
	}

	return actionstypes.Balance{
		Coins: []actionstypes.Coin{
			{
				Denom:  params.BondDenom,
				Amount: totalAmount.String(),
			},
		},
	}, nil
}
