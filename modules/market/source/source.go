package source

import (
	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"
)

type Source interface {
	GetLeases(height int64) ([]markettypes.Lease, error)
}
