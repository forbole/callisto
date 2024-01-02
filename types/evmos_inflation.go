package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	inflationtypes "github.com/evmos/evmos/v15/x/inflation/types"
)

// EvmosInflationParams represents the parameters of the evmos x/inflation module
type EvmosInflationParams struct {
	inflationtypes.Params
	Height int64
}

// NewEvmosInflationParams allows to build a new EvmosInflationParams instance
func NewEvmosInflationParams(params inflationtypes.Params, height int64) *EvmosInflationParams {
	return &EvmosInflationParams{
		Params: params,
		Height: height,
	}
}

// EvmosInflationData represents all the data cli-queried from the evmos x/inflation module
type EvmosInflationData struct {
	CirculatingSupply  sdk.DecCoins
	EpochMintProvision sdk.DecCoins
	InflationRate      sdk.Dec
	InflationPeriod    uint64
	SkippedEpochs      uint64
	Height             int64
}

// NewEvmosInflationData allows to build a new EvmosInflationData instance
func NewEvmosInflationData(
	circulatingSupply sdk.DecCoins, epochMintProvision sdk.DecCoins, inflationRate sdk.Dec,
	inflationPeriod uint64, skippedEpochs uint64, height int64,
) *EvmosInflationData {
	return &EvmosInflationData{
		CirculatingSupply:  circulatingSupply,
		EpochMintProvision: epochMintProvision,
		InflationRate:      inflationRate,
		InflationPeriod:    inflationPeriod,
		SkippedEpochs:      skippedEpochs,
		Height:             height,
	}
}
