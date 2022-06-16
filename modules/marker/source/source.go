package source

import (
	types "github.com/cosmos/cosmos-sdk/codec/types"
)

type Source interface {
	GetAllMarkers(height int64, status string) ([]*types.Any, error)
}
