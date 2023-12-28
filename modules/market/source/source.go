package source

import (
	markettypes "github.com/akash-network/akash-api/go/node/market/v1beta4"
)

type Source interface {
	GetActiveLeases(height int64) ([]markettypes.QueryLeaseResponse, error)
}
