package source

import inflationtypes "github.com/MonikaCat/em-ledger/x/inflation/types"

type Source interface {
	GetInflation(height int64) (inflationtypes.InflationState, error)
}
