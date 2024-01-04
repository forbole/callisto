package remote

import (
	"fmt"

	authoritytypes "github.com/MonikaCat/em-ledger/x/authority/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v4/node/local"

	authoritysource "github.com/forbole/bdjuno/v4/modules/authority/source"
)

var (
	_ authoritysource.Source = &Source{}
)

// Source implements authoritysource.Source using a local node
type Source struct {
	*local.Source
	client authoritytypes.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, client authoritytypes.QueryServer) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// GetMinimumGasPrices implements authoritysource.Source
func (s *Source) GetMinimumGasPrices(height int64) (sdk.DecCoins, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height")
	}

	res, err := s.client.GasPrices(sdk.WrapSDKContext(ctx), &authoritytypes.QueryGasPricesRequest{})
	if err != nil {
		return nil, fmt.Errorf("errror while reading gas prices: %s", err)
	}

	return res.MinGasPrices, nil
}
