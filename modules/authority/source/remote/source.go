package remote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/v2/node/remote"
	authoritytypes "github.com/e-money/em-ledger/x/authority/types"

	authoritysource "github.com/forbole/bdjuno/v2/modules/authority/source"
)

var (
	_ authoritysource.Source = &Source{}
)

// Source implements authoritysource.Source using a remote node
type Source struct {
	*remote.Source
	client authoritytypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, client authoritytypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// GetMinimumGasPrices implements authoritysource.Source
func (s *Source) GetMinimumGasPrices(height int64) (sdk.DecCoins, error) {
	res, err := s.client.GasPrices(s.Ctx, &authoritytypes.QueryGasPricesRequest{}, remote.GetHeightRequestHeader(height))
	if err != nil {
		return nil, fmt.Errorf("errror while querying gas prices: %s", err)
	}

	return res.MinGasPrices, nil
}
