package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v4/node/local"

	liquidstakingtypes "github.com/crescent-network/crescent/v4/x/liquidstaking/types"
	liquidstakingsource "github.com/forbole/bdjuno/v3/modules/liquidstaking/source"
)

var (
	_ liquidstakingsource.Source = &Source{}
)

// Source implements liquidstakingsource.Source using a local node
type Source struct {
	*local.Source
	querier liquidstakingtypes.QueryServer
}

// NewSource returns a new Source instace
func NewSource(source *local.Source, querier liquidstakingtypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements liquidstakingsource.Source
func (s Source) Params(height int64) (liquidstakingtypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return liquidstakingtypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &liquidstakingtypes.QueryParamsRequest{})
	if err != nil {
		return liquidstakingtypes.Params{}, err
	}

	return res.Params, nil
}

// GetLiquidValidators implements liquidstakingtypes.Source
func (s Source) GetLiquidValidators(height int64) ([]liquidstakingtypes.LiquidValidatorState, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return []liquidstakingtypes.LiquidValidatorState{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.LiquidValidators(sdk.WrapSDKContext(ctx), &liquidstakingtypes.QueryLiquidValidatorsRequest{})
	if err != nil {
		return []liquidstakingtypes.LiquidValidatorState{}, err
	}

	return res.LiquidValidators, nil
}

// GetLiquidStakingStates implements liquidstakingtypes.Source
func (s Source) GetLiquidStakingStates(height int64) (liquidstakingtypes.NetAmountState, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return liquidstakingtypes.NetAmountState{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.States(sdk.WrapSDKContext(ctx), &liquidstakingtypes.QueryStatesRequest{})
	if err != nil {
		return liquidstakingtypes.NetAmountState{}, nil
	}

	return res.NetAmountState, nil
}
