package remote

import (
	liquidstakingtypes "github.com/crescent-network/crescent/x/liquidstaking/types"
	"github.com/forbole/juno/v3/node/remote"

	liquidstakingsource "github.com/forbole/bdjuno/v3/modules/liquidstaking/source"
)

var (
	_ liquidstakingsource.Source = &Source{}
)

// Source implements liquidstakingtypes.Source using a remote node
type Source struct {
	*remote.Source
	querier liquidstakingtypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier liquidstakingtypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements liquidstakingtypes.Source
func (s Source) Params(height int64) (liquidstakingtypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &liquidstakingtypes.QueryParamsRequest{})
	if err != nil {
		return liquidstakingtypes.Params{}, nil
	}

	return res.Params, nil
}

// GetLiquidValidators implements liquidstakingtypes.Source
func (s Source) GetLiquidValidators(height int64) ([]liquidstakingtypes.LiquidValidatorState, error) {
	res, err := s.querier.LiquidValidators(remote.GetHeightRequestContext(s.Ctx, height), &liquidstakingtypes.QueryLiquidValidatorsRequest{})
	if err != nil {
		return []liquidstakingtypes.LiquidValidatorState{}, nil
	}

	return res.LiquidValidators, nil
}
