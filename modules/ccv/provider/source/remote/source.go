package remote

import (
	providertypes "github.com/cosmos/interchain-security/v3/x/ccv/provider/types"
	"github.com/forbole/juno/v5/node/remote"

	providersource "github.com/forbole/bdjuno/v4/modules/ccv/provider/source"
)

var (
	_ providersource.Source = &Source{}
)

// Source implements providersource.Source using a remote node
type Source struct {
	*remote.Source
	querier providertypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, querier providertypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetValidatorProviderAddr implements providersource.Source
func (s Source) GetValidatorProviderAddr(height int64, chainID, consumerAddress string) (string, error) {
	res, err := s.querier.QueryValidatorProviderAddr(remote.GetHeightRequestContext(s.Ctx, height), &providertypes.QueryValidatorProviderAddrRequest{ChainId: chainID, ConsumerAddress: consumerAddress})
	if err != nil {
		return "", err
	}

	return res.ProviderAddress, nil
}