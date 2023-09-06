package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	inflationtypes "github.com/evmos/evmos/v14/x/inflation/types"
)

type Source interface {
	Params(height int64) (inflationtypes.Params, error)
	CirculatingSupply(height int64) (sdk.DecCoin, error)
	EpochMintProvision(height int64) (sdk.DecCoin, error)
	InflationRate(height int64) (sdk.Dec, error)
	InflationPeriod(height int64) (uint64, error)
	SkippedEpochs(height int64) (uint64, error)
}
