package local

// import (
// 	"fmt"

// 	markertypes "github.com/MonCatCat/provenance/x/marker/types"
// 	types "github.com/cosmos/cosmos-sdk/codec/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/types/query"
// 	markersource "github.com/forbole/bdjuno/v4/modules/marker/source"
// 	"github.com/forbole/juno/v5/node/local"
// )

// var (
// 	_ markersource.Source = &Source{}
// )

// // Source implements markersource.Source using a local node
// type Source struct {
// 	*local.Source
// 	querier markertypes.QueryServer
// }

// // NewSource returns a new Source instace
// func NewSource(source *local.Source, querier markertypes.QueryServer) *Source {
// 	return &Source{
// 		Source:  source,
// 		querier: querier,
// 	}
// }

// // GetAllMarkers implements markersource.Source
// func (s Source) GetAllMarkers(height int64) ([]*types.Any, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return nil, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	var markers []*types.Any
// 	var nextKey []byte
// 	var stop = false
// 	for !stop {
// 		res, err := s.querier.AllMarkers(
// 			sdk.WrapSDKContext(ctx),
// 			&markertypes.QueryAllMarkersRequest{
// 				Pagination: &query.PageRequest{
// 					Key:   nextKey,
// 					Limit: 100, // Query 100 markers at a time
// 				},
// 			},
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		nextKey = res.Pagination.NextKey
// 		stop = len(res.Pagination.NextKey) == 0
// 		markers = append(markers, res.Markers...)
// 	}

// 	return markers, nil
// }
