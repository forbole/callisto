package source

import (
	markettypes "github.com/akash-network/node/x/market/types/v1beta2"
)

type Source interface {
	GetActiveLeases(height int64) ([]markettypes.QueryLeaseResponse, error)
}
