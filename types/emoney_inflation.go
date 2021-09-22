package types

import (
	"time"

	inflationtypes "github.com/e-money/em-ledger/x/inflation/types"
)

type EmoneyInflation struct {
	Inflation         inflationtypes.InflationAssets
	LastAppliedTime   time.Time
	LastAppliedHeight int64
	Height            int64
}

// NewEmoneyInflation builds a new EmoneyInflation instance
func NewEmoneyInflation(
	inflation inflationtypes.InflationAssets,
	lastAppliedTime time.Time,
	lastAppliedHeight int64,
	height int64,
) EmoneyInflation {
	return EmoneyInflation{
		Inflation:         inflation,
		LastAppliedTime:   lastAppliedTime,
		LastAppliedHeight: lastAppliedHeight,
		Height:            height,
	}
}
