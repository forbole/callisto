package source

import (
	merkettypes "github.com/ovrclk/akash/x/market/types/v1beta2"
)

type Source interface {
	GetLeases(height int64, ownerAddress string) (merkettypes.Leases, error)
}
