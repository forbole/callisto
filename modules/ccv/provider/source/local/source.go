package local

// import (
// 	"fmt"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	ccvprovidertypes "github.com/cosmos/interchain-security/v2/x/ccv/provider/types"
// 	ccvprovidersource "github.com/forbole/bdjuno/v4/modules/ccv/provider/source"
// 	"github.com/forbole/juno/v5/node/local"
// )

// var (
// 	_ ccvprovidersource.Source = &Source{}
// )

// // Source implements ccvprovidersource.Source using a local node
// type Source struct {
// 	*local.Source
// 	querier ccvprovidertypes.QueryServer
// }

// // NewSource implements a new Source instance
// func NewSource(source *local.Source, querier ccvprovidertypes.QueryServer) *Source {
// 	return &Source{
// 		Source:  source,
// 		querier: querier,
// 	}
// }

// // GetAllConsumerChains implements ccvprovidersource.Source
// func (s Source) GetAllConsumerChains(height int64) ([]*ccvprovidertypes.Chain, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return []*ccvprovidertypes.Chain{}, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	res, err := s.querier.QueryConsumerChains(sdk.WrapSDKContext(ctx), &ccvprovidertypes.QueryConsumerChainsRequest{})
// 	if err != nil {
// 		return []*ccvprovidertypes.Chain{}, nil
// 	}

// 	return res.Chains, nil

// }
