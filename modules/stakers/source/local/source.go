package local

// import (
// 	"fmt"

// 	stakersquerytypes "github.com/KYVENetwork/chain/x/query/types"
// 	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/types/query"
// 	"github.com/forbole/juno/v4/node/local"

// 	stakerssource "github.com/forbole/bdjuno/v4/modules/stakers/source"
// )

// var (
// 	_ stakerssource.Source = &Source{}
// )

// // Source implements stakerssource.Source using a local node
// type Source struct {
// 	*local.Source
// 	querier        stakerstypes.QueryServer
// 	stakersQuerier stakersquerytypes.QueryStakersClient
// }

// // NewSource returns a new Source instace
// func NewSource(source *local.Source, querier stakerstypes.QueryServer, stakersQuerier stakersquerytypes.QueryStakersClient) *Source {
// 	return &Source{
// 		Source:  source,
// 		querier: querier,
//		stakersQuerier: stakersQuerier,
// 	}
// }

// // Params implements stakerssource.Source
// func (s Source) Params(height int64) (stakerstypes.Params, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return stakerstypes.Params{}, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &stakerstypes.QueryParamsRequest{})
// 	if err != nil {
// 		return stakerstypes.Params{}, err
// 	}

// 	return res.Params, nil
// }

// // Stakers implements stakerssource.Source
// func (s Source) Stakers(height int64) ([]stakersquerytypes.FullStaker, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return nil, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	var stakers []stakersquerytypes.FullStaker
// 	var nextKey []byte
// 	var stop = false
// 	for !stop {
// 		res, err := s.stakersQuerier.Stakers(
// 			sdk.WrapSDKContext(ctx),
// 			&stakersquerytypes.QueryStakersRequest{
// 				Pagination: &query.PageRequest{
// 					Key:   nextKey,
// 					Limit: 100, // Query 100 stakers at time
// 				},
// 			},
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		nextKey = res.Pagination.NextKey
// 		stop = len(res.Pagination.NextKey) == 0
// 		stakers = append(stakers, res.Stakers...)
// 	}

// 	return stakers, nil
// }
