package types

import (
	"time"

	inflationtypes "github.com/e-money/em-ledger/x/inflation/types"
)

type EmoneyInflation struct {
	InflationAssets   []inflationtypes.InflationAsset
	LastAppliedTime   time.Time
	LastAppliedHeight int64
	Height            int64
}

// NewEmoneyInfaltion allows to build a new EmoneyInflation instance
func NewEmoneyInfaltion(
	state inflationtypes.InflationState, height int64,
) EmoneyInflation {
	return EmoneyInflation{
		InflationAssets:   state.InflationAssets,
		LastAppliedTime:   state.LastAppliedTime,
		LastAppliedHeight: state.LastAppliedHeight.Int64(),
		Height:            height,
	}
}
