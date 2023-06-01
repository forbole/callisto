package local

// import (
// 	"fmt"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	bundlestypes "github.com/KYVENetwork/chain/x/bundles/types"
// 	"github.com/forbole/juno/v4/node/local"

// 	bundlessource "github.com/forbole/bdjuno/v4/modules/bundles/source"
// )

// var (
// 	_ bundlessource.Source = &Source{}
// )

// // Source implements bundlessource.Source using a local node
// type Source struct {
// 	*local.Source
// 	querier bundlestypes.QueryServer
// }

// // NewSource returns a new Source instace
// func NewSource(source *local.Source, querier bundlestypes.QueryServer) *Source {
// 	return &Source{
// 		Source:  source,
// 		querier: querier,
// 	}
// }

// // Params implements bundlessource.Source
// func (s Source) Params(height int64) (bundlestypes.Params, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return bundlestypes.Params{}, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &bundlestypes.QueryParamsRequest{})
// 	if err != nil {
// 		return bundlestypes.Params{}, err
// 	}

// 	return res.Params, nil
// }
