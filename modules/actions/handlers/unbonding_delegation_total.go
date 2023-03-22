package handlers

import (
	"fmt"
	"math/big"

	"github.com/forbole/bdjuno/v4/modules/actions/types"

	"github.com/rs/zerolog/log"
)

func UnbondingDelegationsTotal(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msg("executing unbonding delegation total action")

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	// Get all unbonding delegations for given delegator address
	unbondingDelegations, err := ctx.Sources.StakingSource.GetUnbondingDelegations(height, payload.GetAddress(), nil)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator unbonding delegations: %s", err)
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

	return types.Balance{
		Coins: []types.Coin{
			{
				Denom:  params.BondDenom,
				Amount: totalAmount.String(),
			},
		},
	}, nil
}
