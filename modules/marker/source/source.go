package source

import (
	types "github.com/cosmos/cosmos-sdk/codec/types"
	markertypes "github.com/provenance-io/provenance/x/marker/types"
)

type Source interface {
	GetAllMarkers(height int64) ([]*types.Any, error)
}
