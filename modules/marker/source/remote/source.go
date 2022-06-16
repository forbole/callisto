package remote

import (
	types "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	markersource "github.com/forbole/bdjuno/v3/modules/marker/source"
	"github.com/forbole/juno/v3/node/remote"
	markertypes "github.com/provenance-io/provenance/x/marker/types"
)

var (
	_ markersource.Source = &Source{}
)

// Source implements markersource.Source using a remote node
type Source struct {
	*remote.Source
	querier markertypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier markertypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetAllMarkers implements markersource.Source
func (s Source) GetAllMarkers(height int64) ([]*types.Any, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var markers []*types.Any
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.querier.AllMarkers(
			ctx,
			&markertypes.QueryAllMarkersRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 accounts at a time
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		markers = append(markers, res.Markers...)

	}

	return markers, nil
}
