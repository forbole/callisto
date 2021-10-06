package source

import inflationtypes "github.com/e-money/em-ledger/x/inflation/types"

type Source interface {
	GetInflation(height int64) (inflationtypes.InflationState, error)
}
