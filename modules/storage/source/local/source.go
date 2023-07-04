package local

// import (
// 	"fmt"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/forbole/juno/v4/node/local"

// 	storagesource "github.com/forbole/bdjuno/v4/modules/storage/source"
// 	storagetypes "github.com/jackalLabs/canine-chain/v3/x/storage/types"
// )

// var (
// 	_ storagesource.Source = &Source{}
// )

// // Source implements storagesource.Source using a local node
// type Source struct {
// 	*local.Source
// 	querier storagetypes.QueryServer
// }

// // NewSource returns a new Source instace
// func NewSource(source *local.Source, querier storagetypes.QueryServer) *Source {
// 	return &Source{
// 		Source:  source,
// 		querier: querier,
// 	}
// }

// // Params implements storagesource.Source
// func (s Source) Params(height int64) (storagetypes.Params, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return storagetypes.Params{}, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &storagetypes.QueryParamsRequest{})
// 	if err != nil {
// 		return storagetypes.Params{}, err
// 	}

// 	return res.Params, nil
// }

// // Providers implements storagesource.Source
// func (s Source) Providers(height int64) ([]storagetypes.Providers, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return []storagetypes.Providers{}, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	res, err := s.querier.ProvidersAll(sdk.WrapSDKContext(ctx), &storagetypes.QueryAllProvidersRequest{})
// 	if err != nil {
// 		return []storagetypes.Providers{}, nil
// 	}

// 	return res.Providers, nil
// }
