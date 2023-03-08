package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
	ccvconsumersource "github.com/forbole/bdjuno/v4/modules/ccv/consumer/source"
	"github.com/forbole/juno/v4/node/local"
)

var (
	_ ccvconsumersource.Source = &Source{}
)

// Source implements ccvconsumersource.Source using a local node
type Source struct {
	*local.Source
	querier ccvconsumertypes.QueryServer
}

// NewSource implements a new Source instance
func NewSource(source *local.Source, querier ccvconsumertypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetNextFeeDistribution implements ccvconsumersource.Source
func (s Source) GetNextFeeDistribution(height int64) (*ccvconsumertypes.NextFeeDistributionEstimate, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}
	res, err := s.querier.QueryNextFeeDistribution(sdk.WrapSDKContext(ctx), &ccvconsumertypes.QueryNextFeeDistributionEstimateRequest{})
	if err != nil {
		return nil, err
	}

	return res.Data, nil

}
