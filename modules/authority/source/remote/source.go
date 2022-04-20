package remote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/e-money/em-ledger/x/authority/types"
	"github.com/forbole/juno/v3/node/remote"

	authoritysource "github.com/forbole/bdjuno/v3/modules/authority/source"
	"github.com/forbole/bdjuno/v3/utils"
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
	res, err := s.client.GasPrices(utils.GetHeightRequestContext(s.Ctx, height), &authoritytypes.QueryGasPricesRequest{})
	if err != nil {
		return nil, fmt.Errorf("errror while querying gas prices: %s", err)
	}

	return res.MinGasPrices, nil
}
