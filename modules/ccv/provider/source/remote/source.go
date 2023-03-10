package remote

import (
	ccvprovidertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
	ccvprovidersource "github.com/forbole/bdjuno/v4/modules/ccv/provider/source"
	"github.com/forbole/juno/v4/node/remote"
)

var (
	_ ccvprovidersource.Source = &Source{}
)

// Source implements ccvprovidersource.Source using a remote node
type Source struct {
	*remote.Source
	querier ccvprovidertypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, querier ccvprovidertypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetAllConsumerChains implements ccvprovidersource.Source
func (s Source) GetAllConsumerChains(height int64) ([]*ccvprovidertypes.Chain, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	res, err := s.querier.QueryConsumerChains(ctx, &ccvprovidertypes.QueryConsumerChainsRequest{})
	if err != nil {
		return []*ccvprovidertypes.Chain{}, err
	}

	return res.Chains, nil

}
