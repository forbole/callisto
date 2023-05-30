package local

// import (
// 	"fmt"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
// 	"github.com/forbole/juno/v4/node/local"

// 	stakerssource "github.com/forbole/bdjuno/v4/modules/stakers/source"
// )

// var (
// 	_ stakerssource.Source = &Source{}
// )

// // Source implements stakerssource.Source using a local node
// type Source struct {
// 	*local.Source
// 	querier stakerstypes.QueryServer
// }

// // NewSource returns a new Source instace
// func NewSource(source *local.Source, querier stakerstypes.QueryServer) *Source {
// 	return &Source{
// 		Source:  source,
// 		querier: querier,
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
