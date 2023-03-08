package remote

import (
	ccvconsumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
	ccvconsumersource "github.com/forbole/bdjuno/v4/modules/ccv/consumer/source"
	"github.com/forbole/juno/v4/node/remote"
)

var (
	_ ccvconsumersource.Source = &Source{}
)

// Source implements ccvconsumersource.Source using a remote node
type Source struct {
	*remote.Source
	querier ccvconsumertypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, querier ccvconsumertypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetNextFeeDistribution implements ccvconsumersource.Source
func (s Source) GetNextFeeDistribution(height int64) (*ccvconsumertypes.NextFeeDistributionEstimate, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	res, err := s.querier.QueryNextFeeDistribution(ctx, &ccvconsumertypes.QueryNextFeeDistributionEstimateRequest{})
	if err != nil {
		return nil, err
	}

	return res.Data, nil

}
