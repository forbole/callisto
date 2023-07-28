package local

// import (
// 	"fmt"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/forbole/juno/v4/node/local"

// 	providertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"

// 	providersource "github.com/forbole/bdjuno/v4/modules/ccv/provider/source"
// )

// var (
// 	_ providersource.Source = &Source{}
// )

// // Source implements providersource.Source using a local node
// type Source struct {
// 	*local.Source
// 	querier providertypes.QueryServer
// }

// // NewSource implements a new Source instance
// func NewSource(source *local.Source, querier providertypes.QueryServer) *Source {
// 	return &Source{
// 		Source:  source,
// 		querier: querier,
// 	}
// }

// // GetValidatorProviderAddr implements providersource.Source
// func (s Source) GetValidatorProviderAddr(height int64, chainID, consumerAddress string) (string, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return "", fmt.Errorf("error while loading height: %s", err)
// 	}

// 	res, err := s.querier.QueryValidatorProviderAddr(sdk.WrapSDKContext(ctx), &providertypes.QueryValidatorProviderAddrRequest{})
// 	if err != nil {
// 		return "", nil
// 	}

// 	return res.ProviderAddress, nil
// }
