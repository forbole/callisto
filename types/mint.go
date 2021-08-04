package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MintParams represents the x/mint parameters
type MintParams struct {
	MintDenom           string
	InflationRateChange sdk.Dec
	InflationMax        sdk.Dec
	InflationMin        sdk.Dec
	GoalBonded          sdk.Dec
	BlocksPerYear       uint64
	Height              int64
}

// NewMintParams allows to build a new MintParams instance
func NewMintParams(
	mintDenom string,
	inflationRateChange sdk.Dec,
	inflationMax sdk.Dec,
	inflationMin sdk.Dec,
	goalBonded sdk.Dec,
	blocksPerYear uint64,
	height int64,
) MintParams {
	return MintParams{
		MintDenom:           mintDenom,
		InflationRateChange: inflationRateChange,
		InflationMax:        inflationMax,
		InflationMin:        inflationMin,
		GoalBonded:          goalBonded,
		BlocksPerYear:       blocksPerYear,
		Height:              height,
	}
}
