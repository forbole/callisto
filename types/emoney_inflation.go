package types

import (
	"time"

	inflationtypes "github.com/e-money/em-ledger/x/inflation/types"
)

type EMoneyInflation struct {
	InflationAssets   []inflationtypes.InflationAsset
	LastAppliedTime   time.Time
	LastAppliedHeight int64
	Height            int64
}

// NewEMoneyInfaltion allows to build a new EMoneyInflation instance
func NewEMoneyInfaltion(
	state inflationtypes.InflationState, height int64,
) EMoneyInflation {
	return EMoneyInflation{
		InflationAssets:   state.InflationAssets,
		LastAppliedTime:   state.LastAppliedTime,
		LastAppliedHeight: state.LastAppliedHeight.Int64(),
		Height:            height,
	}
}
